package events

import (
	"errors"

	"grovehub/services/audit/internal/app/validators/events/constants"
)

var ErrInvalidEvent = errors.New(constants.InvalidEventMessage)
