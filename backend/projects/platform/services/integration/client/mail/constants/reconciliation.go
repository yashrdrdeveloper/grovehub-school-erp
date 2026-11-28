package constants

import (
	integrationErrors "grovehub/services/integration/internal/app/constants/error"
	mailConstants "grovehub/services/integration/internal/app/constants/mail"
)

const (
	ReconciliationPath = mailConstants.ReconciliationPath
	ReconciliationScope = mailConstants.ReconciliationScope
	ReconciliationSubject = mailConstants.ReconciliationSubject
)

var ErrInvalidScope = integrationErrors.ErrInvalidScope
