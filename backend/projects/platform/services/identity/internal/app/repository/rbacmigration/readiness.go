package rbacmigration

import (
	"context"
	"errors"
)

var ErrNotReady = errors.New("tenant RBAC migration is not ready")

// Readiness is a necessary migration gate, not authorization or proof that
// every tenant account has its final role configuration.
type Readiness interface {
	Require(ctx context.Context, tenantID string) error
}
