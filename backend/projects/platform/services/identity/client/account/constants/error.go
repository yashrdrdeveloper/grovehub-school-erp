package constants

import accountErrors "grovehub/services/identity/internal/app/constants/error"

// Account errors remain owned by Identity.
var (
	InvalidRequest = accountErrors.InvalidRequest
	Unauthorized = accountErrors.Unauthorized
	Forbidden = accountErrors.Forbidden
	NotFound = accountErrors.NotFound
	Conflict = accountErrors.Conflict
	InvalidState = accountErrors.InvalidState
	Internal = accountErrors.Internal
)
