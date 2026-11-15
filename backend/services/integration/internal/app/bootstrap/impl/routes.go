package impl

import (
    config "grovehub/services/integration/config"
    controller "grovehub/services/integration/internal/app/controller"
    controllerimpl "grovehub/services/integration/internal/app/controller/impl"
    services "grovehub/services/integration/internal/app/services"
)

func (b *builder) newRoutes(connections services.ConnectionService, providers services.ProviderService, push config.GmailPush) controller.Controller {
    return controllerimpl.NewConnectionController(connections,providers,push,b.devAuth)
}
