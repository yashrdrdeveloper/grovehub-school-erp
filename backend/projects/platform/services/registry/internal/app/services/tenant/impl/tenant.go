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
	dto "grovehub/services/registry/internal/app/dto/tenant"
	"grovehub/services/registry/internal/app/entity"
	repository "grovehub/services/registry/internal/app/repository/tenant"
	"grovehub/services/registry/internal/app/services/helpers"
	"grovehub/services/registry/internal/app/services/tenant"
	validator "grovehub/services/registry/internal/app/validators/tenant"
)

// @Service
type TenantService struct {
	repository repository.TenantRepository
}

func NewTenantService(repo repository.TenantRepository) tenant.TenantService {
	return &TenantService{repository: repo}
}

func (s *TenantService) Get(ctx context.Context, projectID, tenantID string) (entity.Tenant, error) {
	if err := validator.Validate(projectID, tenantID); err != nil {
		return entity.Tenant{}, err
	}
	target := scope.Target{Level: scope.Tenant, ProjectID: projectID, TenantID: tenantID}
	if _, err := scope.Require(ctx, target, capability.TenantRead); err != nil {
		return entity.Tenant{}, coreerrors.Wrap(definitions.Forbidden, err)
	}
	result, err := s.repository.Find(ctx, projectID, tenantID)
	return result, helpers.RepositoryError(err)
}

func (s *TenantService) Create(ctx context.Context, projectID string, input dto.CreateInput) (entity.Tenant, error) {
	if err := validator.ValidateCreate(projectID, input); err != nil {
		return entity.Tenant{}, err
	}
	target := scope.Target{Level: scope.Project, ProjectID: projectID}
	if _, err := scope.Require(ctx, target, capability.TenantCreate); err != nil {
		return entity.Tenant{}, coreerrors.Wrap(definitions.Forbidden, err)
	}
	gid, err := id.NewUUIDv7()
	if err != nil {
		return entity.Tenant{}, helpers.RepositoryError(err)
	}
	now := time.Now().UTC()
	value := entity.Tenant{
		TenantID: gid,
		ProjectID: projectID,
		ExternalReference: input.ExternalReference,
		SessionGeneration: validation.InitialGeneration,
		Status: status.Disabled,
		Version: validation.InitialVersion,
		CreatedAt: now,
		UpdatedAt: now,
	}
	result, err := s.repository.Create(ctx, value)
	return result, helpers.RepositoryError(err)
}

func (s *TenantService) ChangeStatus(ctx context.Context, projectID, tenantID string, input dto.StatusInput) (entity.Tenant, error) {
	if err := validator.ValidateStatus(projectID, tenantID, input); err != nil {
		return entity.Tenant{}, err
	}
	target := scope.Target{Level: scope.Tenant, ProjectID: projectID, TenantID: tenantID}
	if _, err := scope.Require(ctx, target, capability.TenantStatus); err != nil {
		return entity.Tenant{}, coreerrors.Wrap(definitions.Forbidden, err)
	}
	result, err := s.repository.ChangeStatus(ctx, projectID, tenantID, input.ExpectedVersion, input.Status)
	return result, helpers.RepositoryError(err)
}
