package authorization

import (
	"context"
	"errors"
)

var (
	ErrScopeRequired = errors.New("authorization snapshot requires tenant and account scope")
	ErrNotFound = errors.New("authorization account not found")
	ErrInvalidSnapshot = errors.New("inconsistent authorization snapshot")
)

// Repository loads authorization state, not a permission decision.
type Repository interface {
	LoadSnapshot(ctx context.Context, tenantID string, accountID string) (Snapshot, error)
}
