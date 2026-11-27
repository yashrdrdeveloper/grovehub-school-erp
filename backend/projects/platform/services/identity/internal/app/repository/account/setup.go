package account

import (
	"context"
	"time"
)

// SetupWriter is separate from admin lifecycle transitions: only the setup
// operation may activate a pending account after consuming its bound grant.
type SetupWriter interface {
	ActivateSetup(ctx context.Context, tenantID string, accountID string, version int64, at time.Time) error
}
