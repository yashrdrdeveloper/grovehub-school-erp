package helpers

import (
	"grovehub/services/identity/internal/app/services/authorization"
	rolerepo "grovehub/services/identity/internal/app/repository/role"
	approvalrepo "grovehub/services/identity/internal/app/repository/approval"
	accountrole "grovehub/services/identity/internal/app/repository/accountrole"
	rolevalidator "grovehub/services/identity/internal/app/validators/role"
	"grovehub/services/identity/internal/app/services/authorization/operations/delegation"
	"errors"
	"strings"

	coreerrors "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors"
	definitions "grovehub/services/identity/internal/app/constants/error"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	sessionrepo "grovehub/services/identity/internal/app/repository/session"
	tokenrepo "grovehub/services/identity/internal/app/repository/passwordtoken"
	auth "grovehub/services/identity/internal/app/services/authentication"
	principal "grovehub/services/identity/internal/app/services/security/principal"
)

// Map translates domain failures into a stable, safe HTTP error envelope.
// Typed errors take precedence; text fallback is limited to legacy operations.
func Map(err error) error {
	if err == nil {
		return nil
	}

	var app *coreerrors.Error
	if errors.As(err, &app) {
		return err
	}

	switch {
	case errors.Is(err, authorization.ErrForbidden):
		return coreerrors.NewError(definitions.Forbidden)
	case errors.Is(err, authorization.ErrUnavailable):
		return coreerrors.NewError(definitions.Unavailable)
	case errors.Is(err, rolerepo.ErrNotFound), errors.Is(err, approvalrepo.ErrNotFound), errors.Is(err, accountrole.ErrNotFound):
		return coreerrors.NewError(definitions.RBACNotFound)
	case errors.Is(err, rolerepo.ErrVersionConflict), errors.Is(err, approvalrepo.ErrConflict), errors.Is(err, accountrole.ErrVersionConflict), errors.Is(err, accountrole.ErrIneligible):
		return coreerrors.NewError(definitions.RBACConflict)
	case errors.Is(err, rolevalidator.ErrInvalidID), errors.Is(err, rolevalidator.ErrInvalidVersion), errors.Is(err, rolevalidator.ErrInvalidSearch), errors.Is(err, approvalrepo.ErrInvalid), errors.Is(err, accountrole.ErrInvalidAssignment), errors.Is(err, delegation.ErrInvalidSet):
		return coreerrors.NewError(definitions.InvalidRequest)
	case errors.Is(err, accountrepo.ErrNotFound):
		return coreerrors.NewError(definitions.NotFound)
	case errors.Is(err, accountrepo.ErrConflict), errors.Is(err, accountrepo.ErrVersionConflict):
		return coreerrors.NewErrorWithMetadata(definitions.Conflict, map[string]any{"resource": "account"})
	case errors.Is(err, tokenrepo.ErrAttemptsExceeded):
		return coreerrors.NewError(definitions.SetupAttemptsExceeded)
	case errors.Is(err, sessionrepo.ErrInvalidSession), errors.Is(err, tokenrepo.ErrInvalidToken):
		return coreerrors.NewError(definitions.InvalidToken)
	case errors.Is(err, auth.ErrInvalidCredentials), errors.Is(err, auth.ErrInvalidRefreshToken):
		return coreerrors.NewError(definitions.InvalidCredentials)
	case errors.Is(err, principal.ErrAuthenticationRequired):
		return coreerrors.NewError(definitions.Unauthorized)
	case errors.Is(err, principal.ErrTenantForbidden), errors.Is(err, principal.ErrAdminForbidden), errors.Is(err, principal.ErrRecentAuthenticationRequired):
		return coreerrors.NewError(definitions.Forbidden)
	}

	message := err.Error()
	switch {
	case strings.Contains(message, "invalid credentials"), strings.Contains(message, "invalid refresh token"):
		return coreerrors.NewError(definitions.InvalidCredentials)
	case strings.Contains(message, "invalid proof"), strings.Contains(message, "invalid reset grant"), strings.Contains(message, "invalid or expired"):
		return coreerrors.NewError(definitions.InvalidToken)
	case strings.Contains(message, "authentication required"):
		return coreerrors.NewError(definitions.Unauthorized)
	case strings.Contains(message, "forbidden"), strings.Contains(message, "unauthorized tenant"), strings.Contains(message, "recent authentication"):
		return coreerrors.NewError(definitions.Forbidden)
	case strings.Contains(message, "already"), strings.Contains(message, "version conflict"), strings.Contains(message, "status transition"):
		return coreerrors.NewError(definitions.Conflict)
	case strings.Contains(message, "username"):
		return coreerrors.NewErrorWithMetadata(definitions.InvalidRequest, map[string]any{"field": "username"})
	case strings.Contains(message, "invalid"), strings.Contains(message, "must be"), strings.Contains(message, "required"):
		return coreerrors.NewError(definitions.InvalidRequest)
	default:
		return coreerrors.NewError(definitions.Internal)
	}
}
