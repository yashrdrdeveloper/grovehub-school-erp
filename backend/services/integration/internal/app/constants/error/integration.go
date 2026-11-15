package errorcode

import (
	"errors"
	coreerrors "grovehub/libs/core/errors"
)

var (
	Invalid   = coreerrors.NewDefinition("INTEGRATION_INVALID", "Invalid integration request")
	Forbidden = coreerrors.NewDefinition("INTEGRATION_FORBIDDEN", "Integration access denied")
	NotFound  = coreerrors.NewDefinition("INTEGRATION_NOT_FOUND", "Integration connection not found")
	NotReady  = coreerrors.NewDefinition("INTEGRATION_NOT_READY", "Provider authorization not implemented")
	Failed    = coreerrors.NewDefinition("INTEGRATION_FAILED", "Integration operation failed")
)

var (
	ErrUnsupportedAction   = errors.New("integration action is unsupported")
	ErrInvalidActionInput  = errors.New("invalid integration action input")
	ErrProviderUnavailable = errors.New("integration provider is not configured")
	ErrUnknownProvider     = errors.New("unknown integration provider")
	ErrDuplicateProvider   = errors.New("duplicate integration provider")
)
