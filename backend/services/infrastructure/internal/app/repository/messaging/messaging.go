package messaging

import (
    "context"
    "time"

    dto "grovehub/services/infrastructure/internal/app/dto/messaging"
    "grovehub/services/infrastructure/internal/app/entity"
)

// Repository owns durable delivery state and lives behind Infrastructure.
type Repository interface {
    Enqueue(ctx context.Context, event dto.Event) error
    Claim(ctx context.Context, limit int, lease time.Duration) ([]entity.Messaging, error)
    MarkSent(ctx context.Context, eventID string) error
    MarkRetry(ctx context.Context, eventID string, attempt int) error
    PruneSent(ctx context.Context, olderThan time.Time, limit int) (int64, error)
}
