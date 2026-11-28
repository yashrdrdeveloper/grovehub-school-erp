package events

import (
	coreerrors "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors"
	"grovehub/services/audit/internal/app/validators/events/constants"
)

var InvalidEvent = coreerrors.NewDefinition(constants.InvalidEventCode, constants.InvalidEventMessage)
var ErrInvalidEvent = coreerrors.NewError(InvalidEvent)
