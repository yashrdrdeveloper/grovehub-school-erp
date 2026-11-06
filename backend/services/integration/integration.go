package integration

import (
 "log/slog"
 "net/http"

 "gorm.io/gorm"

 "grovehub/services/integration/controller"
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
func New(_ *slog.Logger,db *gorm.DB)(Module,error){
 registry,err:=providers.NewRegistry(gmail.New(),outlook.New())
 if err!=nil{return nil,err}
 repository:=repositoryimpl.NewConnectionRepository(db)
 service:=serviceimpl.NewConnectionService(repository,registry)
 return &module{controller:controllerimpl.NewConnectionController(service)},nil
}
func (m *module)RegisterHTTP(mux *http.ServeMux){routes.RegisterHTTP(mux,m.controller)}
