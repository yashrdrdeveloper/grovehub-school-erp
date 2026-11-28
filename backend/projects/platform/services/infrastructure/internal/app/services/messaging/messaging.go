package messaging

import (
    "context"

    dto "grovehub/services/infrastructure/internal/app/dto/messaging"
    inframessaging "grovehub/services/infrastructure/client/messaging"
)

type Service interface {
    Publish(ctx context.Context, topic string, payload []byte, metadata map[string]string) error
    PublishWithOptions(ctx context.Context, topic string, payload []byte, metadata map[string]string, options inframessaging.PublishOptions) error
    Consume(ctx context.Context, topic string, handler dto.MessageHandler) error
    Health(ctx context.Context) error
    Close() error
}
