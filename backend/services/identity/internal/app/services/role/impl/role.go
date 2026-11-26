package impl

import (
	"context"

	"grovehub/services/identity/internal/app/constants/permission"
	roledto "grovehub/services/identity/internal/app/dto/role"
	rolemapper "grovehub/services/identity/internal/app/mappers/role"
	rolerepo "grovehub/services/identity/internal/app/repository/role"
	"grovehub/services/identity/internal/app/services/authorization"
	"grovehub/services/identity/internal/app/services/role"
	"grovehub/services/identity/internal/app/services/security/principal"
	rolevalidator "grovehub/services/identity/internal/app/validators/role"
)

// @Service
type Service struct {
	repository rolerepo.Repository
	authorizer authorization.Service
}

var _ role.Service = (*Service)(nil)

func New(repository rolerepo.Repository, authorizer authorization.Service) role.Service {
	return &Service{repository: repository, authorizer: authorizer}
}

func (service *Service) authorize(ctx context.Context, code permission.Code) (principal.Principal, error) {
	actor, err := principal.RequireTenant(ctx)
	if err != nil {
		return principal.Principal{}, err
	}
	if service.repository == nil || service.authorizer == nil {
		return principal.Principal{}, authorization.ErrUnavailable
	}
	if err := service.authorizer.Require(ctx, code); err != nil {
		return principal.Principal{}, err
	}
	return actor, nil
}

func (service *Service) CreateDraft(ctx context.Context, request roledto.CreateDraft) (roledto.Detail, error) {
	actor, err := service.authorize(ctx, permission.RoleCreate)
	if err != nil {
		return roledto.Detail{}, err
	}
	request, err = rolevalidator.NormalizeCreateDraft(request)
	if err != nil {
		return roledto.Detail{}, err
	}
	detail, err := service.repository.CreateDraft(ctx, rolerepo.Draft{
		TenantID: actor.TenantID, ActorID: actor.AccountID, Name: request.Name,
		Description: request.Description, Personas: request.Personas,
	})
	if err != nil {
		return roledto.Detail{}, err
	}
	return rolemapper.ToDetail(detail), nil
}

func (service *Service) Get(ctx context.Context, roleID string) (roledto.Detail, error) {
	actor, err := service.authorize(ctx, permission.RoleView)
	if err != nil {
		return roledto.Detail{}, err
	}
	roleID, err = rolevalidator.NormalizeID(roleID)
	if err != nil {
		return roledto.Detail{}, err
	}
	detail, err := service.repository.FindByID(ctx, actor.TenantID, roleID)
	if err != nil {
		return roledto.Detail{}, err
	}
	return rolemapper.ToDetail(detail), nil
}

func (service *Service) UpdateMetadata(ctx context.Context, roleID string, request roledto.UpdateMetadata) (roledto.Detail, error) {
	actor, err := service.authorize(ctx, permission.RoleUpdate)
	if err != nil {
		return roledto.Detail{}, err
	}
	roleID, err = rolevalidator.NormalizeID(roleID)
	if err != nil {
		return roledto.Detail{}, err
	}
	request, err = rolevalidator.NormalizeUpdateMetadata(request)
	if err != nil {
		return roledto.Detail{}, err
	}
	detail, err := service.repository.UpdateMetadata(ctx, rolerepo.MetadataUpdate{
		Draft: rolerepo.Draft{
			TenantID: actor.TenantID, ActorID: actor.AccountID, Name: request.Name,
			Description: request.Description, Personas: request.Personas,
		},
		RoleID: roleID, Version: request.Version,
	})
	if err != nil {
		return roledto.Detail{}, err
	}
	return rolemapper.ToDetail(detail), nil
}
