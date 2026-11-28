package messaging

import (
    dto "grovehub/services/infrastructure/internal/app/dto/messaging"
)

// Validator validates durable event commands independently of persistence.
type Validator interface {
    ValidateEvent(event dto.Event) error
}
