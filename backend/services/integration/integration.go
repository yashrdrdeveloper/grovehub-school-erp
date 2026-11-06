package integration

import (
 "log/slog"
 "net/http"

 "gorm.io/gorm"

 "grovehub/services/integration/controller"
 integrationconfig "grovehub/services/integration/config"
 controllerimpl "grovehub/services/integration/controller/impl"
 "grovehub/services/integration/providers"
 "grovehub/services/integration/providers/gmail"
 "grovehub/services/integration/providers/outlook"
 repositoryimpl "grovehub/services/integration/repository/impl"
 "grovehub/services/integration/routes"
 serviceimpl "grovehub/services/integration/services/impl"
)

type Module interface{RegisterHTTP(*http.ServeMux)}
type module struct{controller controller.ConnectionController}
func New(logger *slog.Logger,db *gorm.DB)(Module,error){
 registry,err:=providers.NewRegistry(gmail.New(),outlook.New())
 if err!=nil{return nil,err}
 repository:=repositoryimpl.NewConnectionRepository(db)
 gmailConfig,configErr:=integrationconfig.GmailFromEnvironment()
 var client *gmail.Client
 var key []byte
 if configErr==nil {client=gmail.NewClient(gmailConfig);key=gmailConfig.EncryptionKey} else {logger.Warn("gmail provider not configured; OAuth and tests disabled")}
 service,err:=serviceimpl.NewConnectionService(repository,registry,client,key)
 if err!=nil{return nil,err}
 return &module{controller:controllerimpl.NewConnectionController(service)},nil
}
func (m *module)RegisterHTTP(mux *http.ServeMux){routes.RegisterHTTP(mux,m.controller)}
