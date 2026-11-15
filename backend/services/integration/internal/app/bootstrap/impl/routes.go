package impl

import (
    config "grovehub/services/integration/internal/app/config"
    controller "grovehub/services/integration/internal/app/controller"
    controllerimpl "grovehub/services/integration/internal/app/controller/impl"
    connectionimpl "grovehub/services/integration/internal/app/controller/connection/impl"
    providerimpl "grovehub/services/integration/internal/app/controller/provider/impl"
    connections "grovehub/services/integration/internal/app/services/connection"
    providers "grovehub/services/integration/internal/app/services/provider"
)

// Both resource controllers use the existing generated Integration HTTP routes.
func (b *builder) newRoutes(connectionsService connections.ConnectionService, providersService providers.ProviderService, push config.GmailPush) controller.Controller {
    return controllerimpl.New(
        connectionimpl.NewConnectionController(connectionsService,push,b.devAuth),
        providerimpl.NewProviderController(providersService),
    )
}
