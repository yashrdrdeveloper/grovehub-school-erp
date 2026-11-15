package impl

import (
    config "grovehub/services/integration/config"
    controller "grovehub/services/integration/internal/app/controller"
    controllerimpl "grovehub/services/integration/internal/app/controller/connection/impl"
    connections "grovehub/services/integration/internal/app/services/connection"
    providers "grovehub/services/integration/internal/app/services/provider"
)

func (b *builder) newRoutes(connectionsService connections.ConnectionService, providersService providers.ProviderService, push config.GmailPush) controller.Controller {
    return controllerimpl.NewConnectionController(connectionsService,providersService,push,b.devAuth)
}
