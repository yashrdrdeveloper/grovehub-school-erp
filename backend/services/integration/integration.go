package integration

import (
 "log/slog"
 "net/http"
 "context"
 "time"
 integrationservices "grovehub/services/integration/services"

 "gorm.io/gorm"

 "grovehub/services/integration/controller"
 infrastructureclient "grovehub/services/infrastructure/client"
 integrationconfig "grovehub/services/integration/config"
 controllerimpl "grovehub/services/integration/controller/impl"
 "grovehub/services/integration/providers"
 "grovehub/services/integration/providers/gmail"
 "grovehub/services/integration/providers/outlook"
 repositoryimpl "grovehub/services/integration/repository/impl"
 "grovehub/services/integration/routes"
 serviceimpl "grovehub/services/integration/services/impl"
)

type Module interface{RegisterHTTP(*http.ServeMux);Close()}
type module struct{controller controller.ConnectionController;stop chan struct{}}
func New(logger *slog.Logger,db *gorm.DB,messaging infrastructureclient.MessagingClient)(Module,error){
 registry,err:=providers.NewRegistry(gmail.New(),outlook.New())
 if err!=nil{return nil,err}
 repository:=repositoryimpl.NewConnectionRepository(db)
 gmailConfig,configErr:=integrationconfig.GmailFromEnvironment()
 var client *gmail.Client
 var key []byte
 if configErr==nil {client=gmail.NewClient(gmailConfig);key=gmailConfig.EncryptionKey} else {logger.Warn("gmail provider not configured; OAuth and tests disabled")}
 service,err:=serviceimpl.NewConnectionService(repository,repositoryimpl.NewWatchRepository(db),registry,client,key,integrationconfig.GmailPushFromEnvironment().Topic,messaging)
 if err!=nil{return nil,err}
 push:=integrationconfig.GmailPushFromEnvironment()
 mod:=&module{controller:controllerimpl.NewConnectionController(service,push),stop:make(chan struct{})}
 if push.Enabled()&&client!=nil {go mod.watchRenewal(service)}
 return mod,nil
}
func (m *module)RegisterHTTP(mux *http.ServeMux){routes.RegisterHTTP(mux,m.controller)}

func(m *module)Close(){close(m.stop)}
func(m *module)watchRenewal(service integrationservices.ConnectionService){
 // Distributed duplicate renewal is safe: Gmail watch registration is
 // idempotent and the DB upsert keeps only the latest expiration.
 ticker:=time.NewTicker(6*time.Hour);defer ticker.Stop()
 for {
  select{
  case <-m.stop:return
  case <-ticker.C:
   ctx,cancel:=context.WithTimeout(context.Background(),5*time.Minute)
   _=service.RenewExpiring(ctx)
   cancel()
  }
 }
}
