// Package impl assembles Integration without exposing its repositories/providers.
package impl

import (
    "errors"
    "log/slog"
    "net/http"

    "gorm.io/gorm"

    integrationclient "grovehub/services/integration/client"
    integrationconfig "grovehub/services/integration/config"
    bootstrap "grovehub/services/integration/internal/app/bootstrap"
    controller "grovehub/services/integration/internal/app/controller"
    "grovehub/services/integration/internal/app/routes"
    connectionservices "grovehub/services/integration/internal/app/services"
    messaging "grovehub/services/infrastructure/client/messaging"
)

type builder struct {
    logger *slog.Logger
    db *gorm.DB
    messaging messaging.MessagingClient
    devAuth integrationconfig.DevAuth
}

var _ bootstrap.Bootstrap = (*builder)(nil)

func New(logger *slog.Logger, db *gorm.DB, events messaging.MessagingClient, auth integrationconfig.DevAuth) bootstrap.Bootstrap {
    return &builder{logger:logger,db:db,messaging:events,devAuth:auth}
}

type runtime struct {
    logger *slog.Logger
    controller controller.Controller
    client integrationclient.Client
    localOnly bool
    stop chan struct{}
    done chan struct{}
}

var _ bootstrap.Runtime = (*runtime)(nil)

func (b *builder) Build() (bootstrap.Runtime,error) {
    if b.db==nil {return nil,errors.New("integration bootstrap requires a database")}
    if b.messaging==nil {return nil,errors.New("integration bootstrap requires Infrastructure messaging")}
    if b.logger==nil {b.logger=slog.Default()}
    if b.devAuth.Enabled {b.logger.Warn("Integration development authentication bypass enabled for loopback clients only; never enable in production")}
    providers,err:=b.newProviders()
    if err!=nil{return nil,err}
    connections,providerService,err:=b.newServices(providers)
    if err!=nil{return nil,err}
    controller:=b.newRoutes(connections,providerService,providers.push)
    client:=b.newClient(connections,providerService)
    stop,done:=b.newEvents(connections,providers.push.Enabled()&&providers.gmail!=nil)
    return &runtime{logger:b.logger,controller:controller,client:client,localOnly:b.devAuth.Enabled,stop:stop,done:done},nil
}

func (module *runtime) RegisterHTTP(mux *http.ServeMux) {routes.RegisterHTTP(mux,module.controller)}
func (module *runtime) Client() integrationclient.Client {return module.client}
func (module *runtime) LocalOnly() bool {return module.localOnly}
func (module *runtime) Close() {close(module.stop);<-module.done}

// The connection service is intentionally referenced via its interface.
var _ connectionservices.ConnectionService
