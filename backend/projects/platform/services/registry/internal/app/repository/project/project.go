package project

import (
	"context"

	"grovehub/services/registry/internal/app/entity"
)

// @Repository entity=Project
type ProjectRepository interface {
	Find(ctx context.Context, projectID string) (entity.Project, error)
}
