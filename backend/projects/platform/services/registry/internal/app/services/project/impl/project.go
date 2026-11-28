package impl

import (
	"context"

	coreerrors "grovehub/libs/core/errors"
	"grovehub/libs/platform/scope"
	"grovehub/services/registry/internal/app/constants/capability"
	definitions "grovehub/services/registry/internal/app/constants/error"
	"grovehub/services/registry/internal/app/entity"
	repository "grovehub/services/registry/internal/app/repository/project"
	"grovehub/services/registry/internal/app/services/helpers"
	"grovehub/services/registry/internal/app/services/project"
	validator "grovehub/services/registry/internal/app/validators/project"
)

// @Service
type ProjectService struct {
	repository repository.ProjectRepository
}

func NewProjectService(repo repository.ProjectRepository) project.ProjectService {
	return &ProjectService{repository: repo}
}

func (s *ProjectService) Get(ctx context.Context, projectID string) (entity.Project, error) {
	if err := validator.Validate(projectID); err != nil {
		return entity.Project{}, err
	}

	target := scope.Target{Level: scope.Project, ProjectID: projectID}
	if _, err := scope.Require(ctx, target, capability.ProjectRead); err != nil {
		return entity.Project{}, coreerrors.Wrap(definitions.Forbidden, err)
	}

	result, err := s.repository.Find(ctx, projectID)
	return result, helpers.ReadError(err)
}
