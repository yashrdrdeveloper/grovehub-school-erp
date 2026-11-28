package project

import (
	"context"

	dto "grovehub/services/registry/internal/app/dto/project"
	"grovehub/services/registry/internal/app/entity"
)

type ProjectService interface {
	Get(ctx context.Context, projectID string) (entity.Project, error)
	Create(ctx context.Context, input dto.CreateInput) (entity.Project, error)
	ChangeStatus(ctx context.Context, projectID string, input dto.StatusInput) (entity.Project, error)
}
