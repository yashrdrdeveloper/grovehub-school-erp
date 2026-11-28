package impl

import (
    integrationclient "grovehub/services/integration/client"
    clientimpl "grovehub/services/integration/internal/app/client/connection/impl"
    connections "grovehub/services/integration/internal/app/services/connection"
    providers "grovehub/services/integration/internal/app/services/provider"
)

func (b *builder) newClient(connectionsService connections.ConnectionService, providersService providers.ProviderService) integrationclient.Client {
    return clientimpl.NewClient(connectionsService,providersService)
}
