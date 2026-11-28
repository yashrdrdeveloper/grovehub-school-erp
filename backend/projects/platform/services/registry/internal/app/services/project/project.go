package project

import (
	"context"

	"grovehub/services/registry/internal/app/entity"
)

type ProjectService interface {
	Get(ctx context.Context, projectID string) (entity.Project, error)
}
