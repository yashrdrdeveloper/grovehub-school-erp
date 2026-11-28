package impl

import (
	client "grovehub/services/infrastructure/client/messaging"
	"grovehub/services/infrastructure/internal/app/services/messaging"
)

// Aliases keep generated signatures expressed in the public client contract.
type (
	PublishOptions = client.PublishOptions
	Handler = client.Handler
)

// messagingClient forwards typed operations to the transport service.
type messagingClient struct {
	service messaging.Service
}

func NewMessagingClient(service messaging.Service) client.Client {
	return &messagingClient{service: service}
}
