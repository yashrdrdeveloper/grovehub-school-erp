package key

import (
	"context"

	"grovehub/services/registry/internal/app/entity"
)

// @Repository entity=WorkloadKey
type KeyRepository interface {
	Create(ctx context.Context, value entity.WorkloadKey) error
	Revoke(ctx context.Context, workloadID, keyID string) error
}
