package authority

import (
	"context"
	"errors"
)

var ErrLastAdministrator = errors.New("tenant must retain an effective administrator")
var ErrBootstrapRequired = errors.New("tenant has no effective administrator; controlled bootstrap required")

// AdministratorReader counts active accounts with all recovery capabilities.
type AdministratorReader interface {
	Count(ctx context.Context, tenantID string) (int64, error)
}

// Guard is an internal transaction invariant, not authorization or approval.
// The operation must use the supplied transaction context, affect only this
// tenant, and return every error. Callers must propagate the guard's error.
type Guard interface {
	Run(ctx context.Context, tenantID string, operation func(context.Context) error) error
}
