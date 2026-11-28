package impl

import (
    integrationclient "grovehub/services/integration/client/connection"
    clientimpl "grovehub/services/integration/internal/app/client/connection/impl"
    connections "grovehub/services/integration/internal/app/services/connection"
    providers "grovehub/services/integration/internal/app/services/provider"
)

func (b *builder) newClient(connectionsService connections.ConnectionService, providersService providers.ProviderService) integrationclient.Client {
    return clientimpl.NewConnectionClient(connectionsService,providersService)
}
