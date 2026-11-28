package impl

import (
    providerservice "grovehub/services/integration/internal/app/services/provider"
    integrationclient "grovehub/services/integration/client/connection"
    integrationservices "grovehub/services/integration/internal/app/services/connection"
)

// integrationClient preserves explicit service boundaries.
// Typed forwarders are generated under .generated/client/impl.
type integrationClient struct {
    connections integrationservices.ConnectionService
    providers providerservice.ProviderService
}

func NewClient(connections integrationservices.ConnectionService, providers providerservice.ProviderService) integrationclient.Client {
    return &integrationClient{connections: connections, providers: providers}
}
