package messaging

import (
	"context"

	infrastructureservices "grovehub/services/infrastructure/services"
)

type Strategy interface {
	Publish(ctx context.Context, topic string, payload []byte, metadata map[string]string) error
	Consume(ctx context.Context, topic string, handler infrastructureservices.MessageHandler) error
	Health(ctx context.Context) error
	Close() error
}
