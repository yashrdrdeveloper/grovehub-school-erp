package workload

import (
	"context"

	"grovehub/services/registry/internal/app/entity"
)

// @Repository entity=Workload
type WorkloadRepository interface {
	Find(ctx context.Context, workloadID string) (entity.Workload, error)
}
