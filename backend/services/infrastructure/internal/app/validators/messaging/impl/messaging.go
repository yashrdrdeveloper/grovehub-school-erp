package impl

import (
    errorcode "grovehub/services/infrastructure/internal/app/constants/error"
    dto "grovehub/services/infrastructure/internal/app/dto/messaging"
    validator "grovehub/services/infrastructure/internal/app/validators/messaging"
)

type messagingValidator struct{}

func New() validator.Validator { return &messagingValidator{} }

func (validator *messagingValidator) ValidateEvent(event dto.Event) error {
    if event.EventID == "" || event.Topic == "" || len(event.Payload) == 0 { return errorcode.ErrInvalidMessage }
    return nil
}
