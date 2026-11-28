package messaging

import (
	"context"
	"grovehub/services/infrastructure/client/messaging/contracts"
)

// @Client messagingClient
type MessagingClient interface {
	// @Proxy service.Publish
	Publish(ctx context.Context, topic string, payload []byte, metadata map[string]string) error
	// @Proxy service.PublishWithOptions
	PublishWithOptions(ctx context.Context, topic string, payload []byte, metadata map[string]string, options contracts.PublishOptions) error
	// @Proxy service.Consume
	Consume(ctx context.Context, topic string, handler contracts.Handler) error
	// @Proxy service.Health
	Health(ctx context.Context) error
}
