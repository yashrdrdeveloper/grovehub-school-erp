package account

import (
	"context"
	"time"

	"grovehub/services/identity/internal/app/entity"
)

type MetadataUpdate struct {
	TenantID string
	AccountID string
	Version int64
	Status entity.AccountStatus
	Username string
	SearchDetails string
	UpdatedAt time.Time
}

type MetadataWriter interface {
	UpdateMetadata(ctx context.Context, request MetadataUpdate) error
}
