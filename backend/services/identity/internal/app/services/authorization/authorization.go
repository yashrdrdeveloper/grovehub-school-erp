package authorization

import (
	"context"
	"errors"

	"grovehub/services/identity/internal/app/constants/permission"
)

var (
	ErrForbidden = errors.New("permission forbidden")
	ErrUnavailable = errors.New("authorization unavailable")
)

// Service reads current authority for the verified actor. It does not trust a
// caller-supplied account, tenant, role or cached browser grant.
type Service interface {
	Require(ctx context.Context, requested permission.Code) error
	Effective(ctx context.Context) ([]permission.Code, error)
}
