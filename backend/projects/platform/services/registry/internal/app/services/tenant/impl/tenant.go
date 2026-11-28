package impl

import (
	"context"

	coreerrors "grovehub/libs/core/errors"
	"grovehub/libs/platform/scope"
	"grovehub/services/registry/internal/app/constants/capability"
	definitions "grovehub/services/registry/internal/app/constants/error"
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
	return result, helpers.ReadError(err)
}
