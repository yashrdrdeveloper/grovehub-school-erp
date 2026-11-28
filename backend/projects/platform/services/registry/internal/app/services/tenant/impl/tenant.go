package impl

import (
	"context"

	coreerrors "grovehub/libs/core/errors"
	"grovehub/libs/core/utils"
	"grovehub/libs/platform/scope"
	"grovehub/services/registry/internal/app/constants/capability"
	definitions "grovehub/services/registry/internal/app/constants/error"
	"grovehub/services/registry/internal/app/entity"
	repository "grovehub/services/registry/internal/app/repository/tenant"
	"grovehub/services/registry/internal/app/services/helpers"
	"grovehub/services/registry/internal/app/services/tenant"
)

// @Service
type service struct {
	repository repository.Repository
}

func New(repo repository.Repository) tenant.Service {
	return &service{repository: repo}
}

func (s *service) Get(ctx context.Context, projectID, tenantID string) (entity.Tenant, error) {
	if utils.ValidateUUID(projectID) != nil || utils.ValidateUUID(tenantID) != nil {
		return entity.Tenant{}, coreerrors.NewError(definitions.Invalid)
	}

	target := scope.Target{Level: scope.Tenant, ProjectID: projectID, TenantID: tenantID}
	if _, err := scope.Require(ctx, target, capability.TenantRead); err != nil {
		return entity.Tenant{}, coreerrors.Wrap(definitions.Forbidden, err)
	}

	result, err := s.repository.Find(ctx, projectID, tenantID)
	return result, helpers.ReadError(err)
}
