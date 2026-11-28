package grant

import (
	"context"

	"grovehub/services/registry/internal/app/entity"
)

// @Repository entity=WorkloadGrant
type GrantRepository interface {
	Create(ctx context.Context, value entity.WorkloadGrant) error
	Revoke(ctx context.Context, workloadID, grantID string) error
}
