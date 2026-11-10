package impl

import (
    infrastructureclient "grovehub/services/infrastructure/client"
    infrastructureservices "grovehub/services/infrastructure/services"
)

// messagingClient preserves the public messaging facade.
// Typed forwarding methods are generated from the client interface annotations.
type messagingClient struct {
    service infrastructureservices.MessagingService
}

func NewMessagingClient(service infrastructureservices.MessagingService) infrastructureclient.MessagingClient {
    return &messagingClient{service:service}
}
