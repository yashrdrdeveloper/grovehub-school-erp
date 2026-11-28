package constants

import integrationErrors "grovehub/services/integration/internal/app/constants/error"

var (
	Invalid = integrationErrors.Invalid
	Forbidden = integrationErrors.Forbidden
	NotFound = integrationErrors.NotFound
	NotReady = integrationErrors.NotReady
	Failed = integrationErrors.Failed
	ErrUnsupportedAction = integrationErrors.ErrUnsupportedAction
	ErrInvalidActionInput = integrationErrors.ErrInvalidActionInput
	ErrProviderUnavailable = integrationErrors.ErrProviderUnavailable
	ErrUnknownProvider = integrationErrors.ErrUnknownProvider
	ErrDuplicateProvider = integrationErrors.ErrDuplicateProvider
	ErrMailDeliveryUncertain = integrationErrors.ErrMailDeliveryUncertain
	ErrMailDeliveryConflict = integrationErrors.ErrMailDeliveryConflict
)
