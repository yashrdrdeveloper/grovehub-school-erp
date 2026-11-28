package tenant

import (
	"context"

	"grovehub/services/registry/internal/app/entity"
)

// @Repository entity=Tenant
type TenantRepository interface {
	Find(ctx context.Context, projectID, tenantID string) (entity.Tenant, error)
}
