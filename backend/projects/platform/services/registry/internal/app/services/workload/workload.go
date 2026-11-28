package workload

import (
	"context"

	dto "grovehub/services/registry/internal/app/dto/workload"
	"grovehub/services/registry/internal/app/entity"
)

type WorkloadService interface {
	Get(ctx context.Context, workloadID string) (entity.Workload, error)
	Create(ctx context.Context, input dto.CreateInput) (entity.Workload, error)
	ChangeStatus(ctx context.Context, workloadID string, input dto.StatusInput) (entity.Workload, error)
	RegisterKey(ctx context.Context, workloadID string, input dto.KeyInput) error
	RevokeKey(ctx context.Context, workloadID, keyID string) error
	CreateGrant(ctx context.Context, workloadID string, input dto.GrantInput) (string, error)
	RevokeGrant(ctx context.Context, workloadID, grantID string) error
}
