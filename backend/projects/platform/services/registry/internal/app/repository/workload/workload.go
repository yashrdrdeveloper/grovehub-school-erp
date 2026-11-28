package workload

import (
	"context"

	"grovehub/services/registry/internal/app/constants/status"
	"grovehub/services/registry/internal/app/entity"
)

// @Repository entity=Workload
type WorkloadRepository interface {
	Find(ctx context.Context, workloadID string) (entity.Workload, error)
	Create(ctx context.Context, value entity.Workload) (entity.Workload, error)
	ChangeStatus(ctx context.Context, workloadID string, expectedVersion int64, next status.Status) (entity.Workload, error)
}
