package impl

import (
	"context"

	coreerrors "grovehub/libs/core/errors"
	"grovehub/libs/core/utils"
	"grovehub/libs/platform/scope"
	"grovehub/services/registry/internal/app/constants/capability"
	definitions "grovehub/services/registry/internal/app/constants/error"
	"grovehub/services/registry/internal/app/entity"
	repository "grovehub/services/registry/internal/app/repository/project"
	"grovehub/services/registry/internal/app/services/helpers"
	"grovehub/services/registry/internal/app/services/project"
)

// @Service
type service struct {
	repository repository.Repository
}

func New(repo repository.Repository) project.Service {
	return &service{repository: repo}
}

func (s *service) Get(ctx context.Context, projectID string) (entity.Project, error) {
	if utils.ValidateUUID(projectID) != nil {
		return entity.Project{}, coreerrors.NewError(definitions.Invalid)
	}

	target := scope.Target{Level: scope.Project, ProjectID: projectID}
	if _, err := scope.Require(ctx, target, capability.ProjectRead); err != nil {
		return entity.Project{}, coreerrors.Wrap(definitions.Forbidden, err)
	}

	result, err := s.repository.Find(ctx, projectID)
	return result, helpers.ReadError(err)
}
