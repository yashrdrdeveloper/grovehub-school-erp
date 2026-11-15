package strategies

import (
    "context"
    dto "grovehub/services/infrastructure/internal/app/dto/messaging"
)

// @Strategy
type Strategy interface {
    Execute(ctx context.Context, topic string, handler dto.MessageHandler) error
}
