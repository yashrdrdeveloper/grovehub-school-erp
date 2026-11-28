package tenant

import (
	"context"

	"grovehub/services/registry/internal/app/constants/status"
	"grovehub/services/registry/internal/app/entity"
)

// @Repository entity=Tenant
type TenantRepository interface {
	Find(ctx context.Context, projectID, tenantID string) (entity.Tenant, error)
	Create(ctx context.Context, value entity.Tenant) (entity.Tenant, error)
	ChangeStatus(ctx context.Context, projectID, tenantID string, expectedVersion int64, next status.Status) (entity.Tenant, error)
}
