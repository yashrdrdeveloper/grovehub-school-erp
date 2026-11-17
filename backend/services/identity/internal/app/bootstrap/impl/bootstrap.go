package impl
import(
 "context"
 "errors"
 "log/slog"
 "net/http"
 "sync"
 "time"
 "gorm.io/gorm"
    coresearch "grovehub/libs/core/search"
 "grovehub/services/identity/internal/app/config"
 configimpl "grovehub/services/identity/internal/app/config/impl"
 bootstrap "grovehub/services/identity/internal/app/bootstrap"
 accountclient "grovehub/services/identity/client/account"
 accountprotected "grovehub/services/identity/internal/app/controller/account/protected"
 accountpublic "grovehub/services/identity/internal/app/controller/account/public"
 authprotected "grovehub/services/identity/internal/app/controller/authentication/protected"
 authpublic "grovehub/services/identity/internal/app/controller/authentication/public"
 "grovehub/services/identity/internal/app/middleware"
 accountservice "grovehub/services/identity/internal/app/services/account"
 invitation "grovehub/services/identity/internal/app/services/registration/operations/invitation"
 "grovehub/services/identity/internal/app/routes"
 foundationclient "grovehub/services/foundation/client/person"
 integrationclient "grovehub/services/integration/client/connection"
)
type builder struct{logger *slog.Logger;db *gorm.DB;foundation foundationclient.Client;integration integrationclient.Client}
func New(logger *slog.Logger,db *gorm.DB,foundation foundationclient.Client,integration integrationclient.Client)bootstrap.Bootstrap{return &builder{logger:logger,db:db,foundation:foundation,integration:integration}}
type runtime struct{
 logger *slog.Logger
 accountPublic accountpublic.AccountController
 accountProtected accountprotected.AccountController
 authPublic authpublic.AuthenticationController
 authProtected authprotected.AuthenticationController
 middleware middleware.Authentication
 accountService accountservice.Service
 worker *invitation.Worker
 workerConfig config.Config
 stop context.CancelFunc
 done chan struct{}
 once sync.Once
}
var _ bootstrap.Runtime=(*runtime)(nil)
func(b *builder)Build()(bootstrap.Runtime,error){
 if b.db==nil||b.foundation==nil{return nil,errors.New("identity requires PostgreSQL and Foundation client")}
 if b.logger==nil{b.logger=slog.Default()}
 c,err:=configimpl.Load();if err!=nil{return nil,err}
 reps:=b.newRepository()
 services:=b.newServices(reps,c)
 authpub,authpriv,accpub,accpriv:=b.newControllers(services)
 m:=middleware.Authentication{Accounts:reps.accounts,Sessions:reps.sessions,AdminIDs:c.AdminIDs,AllowedOrigin:c.AllowedOrigin}
 stopCtx,stop:=context.WithCancel(context.Background())
 r:=&runtime{logger:b.logger,accountPublic:accpub,accountProtected:accpriv,authPublic:authpub,authProtected:authpriv,middleware:m,accountService:services.account,worker:services.worker,workerConfig:c,stop:stop,done:make(chan struct{})}
 go r.runWorker(stopCtx)
 return r,nil
}
func(r *runtime)RegisterHTTP(mux *http.ServeMux){routes.RegisterHTTP(mux,r.authPublic,r.authProtected,r.accountPublic,r.accountProtected,r.middleware)}
func(r *runtime)AccountClient()accountclient.Client{return r}
func(r *runtime)Get(ctx context.Context,id string)(accountclient.Account,error){return r.accountService.Get(ctx,id)}
func(r *runtime)Search(ctx context.Context,q coresearch.Request)(coresearch.Result[accountclient.Account],error){return r.accountService.Search(ctx,q)}
func(r *runtime)Close(){r.once.Do(func(){r.stop();<-r.done})}
func(r *runtime)runWorker(ctx context.Context){
 defer close(r.done)
 if r.worker==nil{return}
 ticker:=time.NewTicker(10*time.Second);defer ticker.Stop()
 for{
  select{case <-ctx.Done():return;case <-ticker.C:
   child,cancel:=context.WithTimeout(ctx,20*time.Second)
   _,err:=r.worker.RunOnce(child);cancel()
   if err!=nil{r.logger.Warn("identity mail delivery retry scheduled")}
  }
 }
}
