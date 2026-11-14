package messaging

import (
    "context"

    dto "grovehub/services/infrastructure/internal/app/dto/messaging"
)

type Service interface {
    Publish(ctx context.Context, topic string, payload []byte, metadata map[string]string) error
    Consume(ctx context.Context, topic string, handler dto.MessageHandler) error
    Health(ctx context.Context) error
    Close() error
}
