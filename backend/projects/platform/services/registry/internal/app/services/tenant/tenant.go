package tenant

import (
	"context"

	dto "grovehub/services/registry/internal/app/dto/tenant"
	"grovehub/services/registry/internal/app/entity"
)

type TenantService interface {
	Get(ctx context.Context, projectID, tenantID string) (entity.Tenant, error)
	Create(ctx context.Context, projectID string, input dto.CreateInput) (entity.Tenant, error)
	ChangeStatus(ctx context.Context, projectID, tenantID string, input dto.StatusInput) (entity.Tenant, error)
}
