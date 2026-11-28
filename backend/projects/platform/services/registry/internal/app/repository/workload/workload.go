package workload

import (
	"context"

	"grovehub/services/registry/internal/app/entity"
)

// @Repository entity=Workload
type Repository interface {
	Find(ctx context.Context, workloadID string) (entity.Workload, error)
}
