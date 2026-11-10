package impl

import (
    integrationclient "grovehub/services/integration/client"
    integrationservices "grovehub/services/integration/services"
)

// integrationClient preserves explicit service boundaries.
// Typed forwarders are generated under .generated/client/impl.
type integrationClient struct {
    connections integrationservices.ConnectionService
    providers integrationservices.ProviderService
}

func NewClient(connections integrationservices.ConnectionService, providers integrationservices.ProviderService) integrationclient.Client {
    return &integrationClient{connections: connections, providers: providers}
}
