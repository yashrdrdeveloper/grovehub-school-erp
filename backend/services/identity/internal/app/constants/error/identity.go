package errorcode

import (
	"net/http"

	coreerrors "grovehub/libs/core/errors"
)

var (
	InvalidRequest = coreerrors.NewDefinition("IDENTITY_REQUEST_INVALID", "Invalid request", http.StatusBadRequest)
	InvalidCredentials = coreerrors.NewDefinition("IDENTITY_INVALID_CREDENTIALS", "Invalid credentials", http.StatusUnauthorized)
	Unauthorized = coreerrors.NewDefinition("IDENTITY_UNAUTHORIZED", "Authentication required", http.StatusUnauthorized)
	Forbidden = coreerrors.NewDefinition("IDENTITY_FORBIDDEN", "Access denied", http.StatusForbidden)
	NotFound = coreerrors.NewDefinition("IDENTITY_ACCOUNT_NOT_FOUND", "Account not found", http.StatusNotFound)
	Conflict = coreerrors.NewDefinition("IDENTITY_ACCOUNT_CONFLICT", "Account already exists or has changed", http.StatusConflict)
	InvalidToken = coreerrors.NewDefinition("IDENTITY_TOKEN_INVALID", "Token is invalid or expired", http.StatusUnauthorized)
	InvalidState = coreerrors.NewDefinition("IDENTITY_ACCOUNT_STATE_INVALID", "Account state does not permit this operation", http.StatusConflict)
	SetupAttemptsExceeded = coreerrors.NewDefinition("IDENTITY_SETUP_ATTEMPTS_EXCEEDED", "Account activation attempts exhausted", http.StatusTooManyRequests)
	RateLimited = coreerrors.NewDefinition("IDENTITY_RATE_LIMITED", "Too many requests", http.StatusTooManyRequests)
	Unavailable = coreerrors.NewDefinition("IDENTITY_UNAVAILABLE", "Authentication service is unavailable", http.StatusServiceUnavailable)
	Internal = coreerrors.NewDefinition("IDENTITY_INTERNAL", "Identity operation failed", http.StatusInternalServerError)
)
