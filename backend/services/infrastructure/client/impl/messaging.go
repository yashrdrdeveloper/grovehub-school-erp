package impl

import (
    infrastructureclient "grovehub/services/infrastructure/client"
    infrastructureservices "grovehub/services/infrastructure/services"
)

// messagingClient preserves the public messaging facade.
// @Client/@Proxy generates the technical forwarding methods.
type messagingClient struct {
    service infrastructureservices.MessagingService
}

func NewMessagingClient(service infrastructureservices.MessagingService) infrastructureclient.MessagingClient {
    return &messagingClient{service:service}
}
