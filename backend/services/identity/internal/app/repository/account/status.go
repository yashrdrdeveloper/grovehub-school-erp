package account

import (
	"context"
	"time"

	"grovehub/services/identity/internal/app/entity"
)

type StatusChange struct {
	TenantID string
	AccountID string
	Version int64
	Before entity.AccountStatus
	After entity.AccountStatus
	UpdatedAt time.Time
}

// StatusWriter intentionally cannot change legacy or normalized role assignments.
type StatusWriter interface {
	ChangeStatus(ctx context.Context, request StatusChange) error
}
