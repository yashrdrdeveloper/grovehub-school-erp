package impl

import (
    integrationclient "grovehub/services/integration/client"
    clientimpl "grovehub/services/integration/internal/app/client/connections/impl"
    services "grovehub/services/integration/internal/app/services"
)

func (b *builder) newClient(connections services.ConnectionService, providers services.ProviderService) integrationclient.Client {
    return clientimpl.NewClient(connections,providers)
}
