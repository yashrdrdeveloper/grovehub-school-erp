package impl

import (
	"context"
	"time"

	coreerrors "grovehub/libs/core/errors"
	"grovehub/libs/core/id"
	"grovehub/libs/platform/scope"
	"grovehub/services/registry/internal/app/constants/capability"
	definitions "grovehub/services/registry/internal/app/constants/error"
	"grovehub/services/registry/internal/app/constants/status"
	"grovehub/services/registry/internal/app/constants/validation"
	dto "grovehub/services/registry/internal/app/dto/project"
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
	return result, helpers.RepositoryError(err)
}

func (s *ProjectService) Create(ctx context.Context, input dto.CreateInput) (entity.Project, error) {
	if err := validator.ValidateCreate(input); err != nil {
		return entity.Project{}, err
	}
	target := scope.Target{Level: scope.Platform}
	if _, err := scope.Require(ctx, target, capability.ProjectCreate); err != nil {
		return entity.Project{}, coreerrors.Wrap(definitions.Forbidden, err)
	}
	gid, err := id.NewUUIDv7()
	if err != nil {
		return entity.Project{}, helpers.RepositoryError(err)
	}
	now := time.Now().UTC()
	value := entity.Project{
		ProjectID: gid,
		Code: input.Code,
		Name: input.Name,
		Status: status.Disabled,
		Version: validation.InitialVersion,
		CreatedAt: now,
		UpdatedAt: now,
	}
	result, err := s.repository.Create(ctx, value)
	return result, helpers.RepositoryError(err)
}

func (s *ProjectService) ChangeStatus(ctx context.Context, projectID string, input dto.StatusInput) (entity.Project, error) {
	if err := validator.ValidateStatus(projectID, input); err != nil {
		return entity.Project{}, err
	}
	target := scope.Target{Level: scope.Project, ProjectID: projectID}
	if _, err := scope.Require(ctx, target, capability.ProjectStatus); err != nil {
		return entity.Project{}, coreerrors.Wrap(definitions.Forbidden, err)
	}
	result, err := s.repository.ChangeStatus(ctx, projectID, input.ExpectedVersion, input.Status)
	return result, helpers.RepositoryError(err)
}
