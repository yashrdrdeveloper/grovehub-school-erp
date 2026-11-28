package tenant

import (
	"context"

	"grovehub/services/registry/internal/app/entity"
)

type Service interface {
	Get(ctx context.Context, projectID, tenantID string) (entity.Tenant, error)
}
