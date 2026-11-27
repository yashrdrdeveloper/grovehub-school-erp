package strategies

import (
	"context"

	messagingdto "grovehub/services/infrastructure/internal/app/dto/messaging"
)

type Strategy interface {
	Publish(ctx context.Context, topic string, payload []byte, metadata map[string]string) error
	Consume(ctx context.Context, topic string, handler messagingdto.MessageHandler) error
	Health(ctx context.Context) error
	Close() error
}
