package project

import (
	"context"

	"grovehub/services/registry/internal/app/constants/status"
	"grovehub/services/registry/internal/app/entity"
)

// @Repository entity=Project
type ProjectRepository interface {
	Find(ctx context.Context, projectID string) (entity.Project, error)
	Create(ctx context.Context, value entity.Project) (entity.Project, error)
	ChangeStatus(ctx context.Context, projectID string, expectedVersion int64, next status.Status) (entity.Project, error)
}
