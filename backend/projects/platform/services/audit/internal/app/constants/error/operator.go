package errorcode

import (
	"net/http"

	coreerrors "grovehub/libs/core/errors"
)

var (
	OperatorInvalid        = coreerrors.NewDefinition("AUDIT_OPERATOR_INVALID", "Invalid operator audit request", http.StatusBadRequest)
	OperatorForbidden      = coreerrors.NewDefinition("AUDIT_OPERATOR_FORBIDDEN", "Access denied", http.StatusForbidden)
	OperatorConflict       = coreerrors.NewDefinition("AUDIT_OPERATOR_CONFLICT", "operator occurrence conflict", http.StatusConflict)
	OperatorUnavailable    = coreerrors.NewDefinition("AUDIT_OPERATOR_UNAVAILABLE", "Audit storage unavailable", http.StatusServiceUnavailable)
	SchoolActivityNotFound = coreerrors.NewDefinition("AUDIT_SCHOOL_ACTIVITY_NOT_FOUND", "School activity not found", http.StatusNotFound)
)
