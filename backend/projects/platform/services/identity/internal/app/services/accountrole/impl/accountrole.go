package impl

import (
	"context"
	"slices"

	"grovehub/services/identity/internal/app/constants/permission"
	dto "grovehub/services/identity/internal/app/dto/accountrole"
	repo "grovehub/services/identity/internal/app/repository/accountrole"
	"grovehub/services/identity/internal/app/services/accountrole"
	"grovehub/services/identity/internal/app/services/accountrole/operations/preview"
	"grovehub/services/identity/internal/app/services/authorization"
	"grovehub/services/identity/internal/app/services/security/principal"
	validator "grovehub/services/identity/internal/app/validators/accountrole"
)

// @Service
type Service struct {
	repository repo.Repository
	authorizer authorization.Service
}

var _ accountrole.Service = (*Service)(nil)

func New(repository repo.Repository, authorizer authorization.Service) accountrole.Service {
	return &Service{repository: repository, authorizer: authorizer}
}

func (service *Service) Get(ctx context.Context, accountID string) (repo.Assignment, error) {
	actor, err := principal.RequireTenant(ctx)
	if err != nil {
		return repo.Assignment{}, err
	}
	if service.repository == nil || service.authorizer == nil {
		return repo.Assignment{}, authorization.ErrUnavailable
	}
	if err := service.authorizer.Require(ctx, permission.AccountRoleView); err != nil {
		return repo.Assignment{}, err
	}
	id, err := validator.NormalizeID(accountID)
	if err != nil {
		return repo.Assignment{}, err
	}
	return service.repository.Find(ctx, actor.TenantID, id)
}

func (service *Service) Preview(ctx context.Context, accountID string, request dto.Replace) (dto.Preview, error) {
	actor, err := principal.RequireTenant(ctx)
	if err != nil {
		return dto.Preview{}, err
	}
	if service.repository == nil || service.authorizer == nil {
		return dto.Preview{}, authorization.ErrUnavailable
	}
	effective, err := service.authorizer.Effective(ctx)
	if err != nil {
		return dto.Preview{}, err
	}
	if !slices.Contains(effective, permission.AccountRoleAssign) {
		return dto.Preview{}, authorization.ErrForbidden
	}
	id, err := validator.NormalizeID(accountID)
	if err != nil {
		return dto.Preview{}, err
	}
	request, err = validator.NormalizeReplace(request)
	if err != nil {
		return dto.Preview{}, err
	}
	current, err := service.repository.Find(ctx, actor.TenantID, id)
	if err != nil {
		return dto.Preview{}, err
	}
	if current.Version != request.Version {
		return dto.Preview{}, repo.ErrVersionConflict
	}
	candidates, err := service.repository.Candidates(ctx, actor.TenantID, preview.RoleIDs(current, request))
	if err != nil {
		return dto.Preview{}, err
	}
	return preview.Evaluate(current, request, candidates, effective)
}
