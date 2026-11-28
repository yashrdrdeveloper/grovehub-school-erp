package impl

import (
 "context"
 "fmt"
 "log/slog"
 "net/http"
 "os"
 "os/signal"
 "sync/atomic"
 "syscall"
 "time"

 "grovehub/launcher/internal/src/bootstrap"
 "grovehub/launcher/internal/src/bootstrap/helpers"
 "grovehub/launcher/internal/src/config"
 serverimpl "grovehub/launcher/internal/src/server/impl"
)

type builder struct {
 logger *slog.Logger
 config config.Config
}
var _ bootstrap.Bootstrap = (*builder)(nil)

func New(logger *slog.Logger) bootstrap.Bootstrap {
 if logger==nil {logger=slog.Default()}
 return &builder{logger:logger,config:config.FromEnv()}
}

func(b *builder) Build()(result bootstrap.Runtime,err error){
 cleanup:=&helpers.Cleanup{}
 ctx,stop:=signal.NotifyContext(context.Background(),os.Interrupt,syscall.SIGTERM)
 cleanup.Add(stop)
 defer func(){if err!=nil{cleanup.Close()}}()

 tenantState,tenantCheck,err:=b.newTenantState()
 if err!=nil{return nil,err}

 databaseCtx,cancelDatabase:=context.WithTimeout(ctx,5*time.Second)
 db,err:=NewPostgresDB(databaseCtx)
 cancelDatabase()
 if err!=nil{return nil,fmt.Errorf("bootstrap PostgreSQL: %w",err)}
 sqlDB,err:=db.DB()
 if err!=nil{return nil,fmt.Errorf("access PostgreSQL connection: %w",err)}
 cleanup.Add(func(){_=sqlDB.Close()})

 modules,err:=b.newServices(ctx,db,tenantState,cleanup)
 if err!=nil{return nil,err}
 mux:=http.NewServeMux()
 draining:=&atomic.Bool{}
 if err:=registerRoutes(mux,db,sqlDB,modules,draining);err!=nil{return nil,err}
 b.startWorkers(ctx,db,tenantState,cleanup)
 gateway,err:=b.newGateway(mux,modules,tenantCheck)
 if err!=nil{return nil,err}
 handler:=newEdgeHandler(b.logger,mux,gateway.External())
 return &runtime{ctx:ctx,server:serverimpl.New(b.config.Address,handler,b.logger,draining),cleanup:cleanup},nil
}
