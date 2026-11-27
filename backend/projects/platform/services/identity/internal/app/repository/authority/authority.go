package authority

import (
	"context"
	"errors"
)

var ErrScopeRequired = errors.New("authority transaction requires tenant scope")

// Locker serializes authority mutations for a tenant. All future role activation,
// assignment, grant and entitlement mutation paths must share this protocol.
type Locker interface {
	Run(ctx context.Context, tenantID string, operation func(context.Context) error) error
}
