package messaging

import (
    "time"

    dto "grovehub/services/infrastructure/internal/app/dto/messaging"
    "grovehub/services/infrastructure/internal/app/entity"
)

// Mapper converts immutable event commands to durable Messaging entities.
type Mapper interface {
    ToEntity(event dto.Event, now time.Time) (entity.Messaging, error)
}
