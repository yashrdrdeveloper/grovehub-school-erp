package events

import (
	coreerrors "grovehub/libs/core/errors"
	"grovehub/services/audit/internal/app/validators/events/constants"
)

var InvalidEvent = coreerrors.NewDefinition(constants.InvalidEventCode, constants.InvalidEventMessage)
var ErrInvalidEvent = coreerrors.NewError(InvalidEvent)
