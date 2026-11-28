package errorcode

import (
	"net/http"

	coreerrors "grovehub/libs/core/errors"
)

var (
	ApprovalExpired = coreerrors.NewDefinition("IDENTITY_APPROVAL_EXPIRED", "Approval request expired", http.StatusConflict)
	ApprovalPolicyUnavailable = coreerrors.NewDefinition("IDENTITY_APPROVAL_POLICY_UNAVAILABLE", "No eligible approval policy is configured", http.StatusConflict)
	ApprovalStateConflict = coreerrors.NewDefinition("IDENTITY_APPROVAL_STATE_CONFLICT", "Approval state or version does not permit this operation", http.StatusConflict)
	ApprovalLastAdministrator = coreerrors.NewDefinition("IDENTITY_APPROVAL_LAST_ADMINISTRATOR", "At least one administrator must remain", http.StatusConflict)
	RBACNotFound = coreerrors.NewDefinition("IDENTITY_RBAC_NOT_FOUND", "Resource not found", http.StatusNotFound)
	RBACConflict = coreerrors.NewDefinition("IDENTITY_RBAC_CONFLICT", "Resource has changed or is not eligible", http.StatusConflict)
	InvalidRequest = coreerrors.NewDefinition("IDENTITY_REQUEST_INVALID", "Invalid request", http.StatusBadRequest)
	InvalidCredentials = coreerrors.NewDefinition("IDENTITY_INVALID_CREDENTIALS", "Invalid credentials", http.StatusUnauthorized)
	Unauthorized = coreerrors.NewDefinition("IDENTITY_UNAUTHORIZED", "Authentication required", http.StatusUnauthorized)
	Forbidden = coreerrors.NewDefinition("IDENTITY_FORBIDDEN", "Access denied", http.StatusForbidden)
	NotFound = coreerrors.NewDefinition("IDENTITY_ACCOUNT_NOT_FOUND", "Account not found", http.StatusNotFound)
	Conflict = coreerrors.NewDefinition("IDENTITY_ACCOUNT_CONFLICT", "Account already exists or has changed", http.StatusConflict)
	InvalidToken = coreerrors.NewDefinition("IDENTITY_TOKEN_INVALID", "Token is invalid or expired", http.StatusUnauthorized)
	InvalidState = coreerrors.NewDefinition("IDENTITY_ACCOUNT_STATE_INVALID", "Account state does not permit this operation", http.StatusConflict)
	ResetAttemptsExceeded = coreerrors.NewDefinition("IDENTITY_RESET_ATTEMPTS_EXCEEDED", "Password reset attempts exhausted. Request a new reset link.", http.StatusTooManyRequests)
	SetupAttemptsExceeded = coreerrors.NewDefinition("IDENTITY_SETUP_ATTEMPTS_EXCEEDED", "Account activation attempts exhausted", http.StatusTooManyRequests)
	RateLimited = coreerrors.NewDefinition("IDENTITY_RATE_LIMITED", "Too many requests", http.StatusTooManyRequests)
	Unavailable = coreerrors.NewDefinition("IDENTITY_UNAVAILABLE", "Authentication service is unavailable", http.StatusServiceUnavailable)
	Internal = coreerrors.NewDefinition("IDENTITY_INTERNAL", "Identity operation failed", http.StatusInternalServerError)
)
