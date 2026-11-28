package tenant

import (
	"context"

	"grovehub/services/registry/internal/app/entity"
)

type TenantService interface {
	Get(ctx context.Context, projectID, tenantID string) (entity.Tenant, error)
}
