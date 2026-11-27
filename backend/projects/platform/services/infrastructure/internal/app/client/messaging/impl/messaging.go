package impl

import (
    infrastructureclient "grovehub/services/infrastructure/client/messaging"
    infrastructureservices "grovehub/services/infrastructure/internal/app/services/messaging"
)

// messagingClient preserves the public messaging facade.
// Typed forwarding methods are generated from the client interface annotations.
type messagingClient struct {
    service infrastructureservices.Service
}

func NewMessagingClient(service infrastructureservices.Service) infrastructureclient.MessagingClient {
    return &messagingClient{service:service}
}
