package impl

import (
	"context"

	dto "grovehub/services/identity/internal/app/dto/role"
	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/constants/persona"
	"grovehub/services/identity/internal/app/events/producers/resource/operations/changed"
	"grovehub/services/identity/internal/app/events/producers/rolemetadata"
	"grovehub/services/identity/internal/app/repository/authority"
	repo "grovehub/services/identity/internal/app/repository/role"
	"grovehub/services/identity/internal/app/services/authorization"
	"grovehub/services/identity/internal/app/services/role"
	"grovehub/services/identity/internal/app/services/security/principal"
	validator "grovehub/services/identity/internal/app/validators/role"
)

// AuditedService exposes only metadata writes, never grants or activation.
type AuditedService struct {
	role.Service
	repository repo.Repository
	authorizer authorization.Service
	locker authority.Locker
	inspector repo.DraftInspector
	recorder rolemetadata.Recorder
}

var _ role.Service = (*AuditedService)(nil)

func NewAudited(repository repo.Repository, authorizer authorization.Service, locker authority.Locker, inspector repo.DraftInspector, recorder rolemetadata.Recorder) role.Service {
	return &AuditedService{Service: New(repository, authorizer), repository: repository, authorizer: authorizer, locker: locker, inspector: inspector, recorder: recorder}
}

func (service *AuditedService) write(ctx context.Context, code permission.Code, operation func(context.Context, principal.Principal) error) error {
	actor, err := principal.RequireTenant(ctx)
	if err != nil {
		return err
	}
	if service.repository == nil || service.authorizer == nil || service.locker == nil || service.inspector == nil || service.recorder == nil {
		return authorization.ErrUnavailable
	}
	return service.locker.Run(ctx, actor.TenantID, func(tx context.Context) error {
		if err := service.authorizer.Require(tx, code); err != nil {
			return err
		}
		return operation(tx, actor)
	})
}

func (service *AuditedService) CreateDraft(ctx context.Context, request dto.CreateDraft) (dto.Detail, error) {
	var result dto.Detail
	err := service.write(ctx, permission.RoleCreate, func(tx context.Context, actor principal.Principal) error {
		var err error
		result, err = service.Service.CreateDraft(tx, request)
		if err != nil {
			return err
		}
		after, err := service.snapshot(tx, actor.TenantID, result.RoleID)
		if err != nil {
			return err
		}
		return service.recorder.Record(tx, nil, after)
	})
	if err != nil {
		return dto.Detail{}, err
	}
	return result, nil
}

func (service *AuditedService) UpdateMetadata(ctx context.Context, roleID string, request dto.UpdateMetadata) (dto.Detail, error) {
	var result dto.Detail
	err := service.write(ctx, permission.RoleUpdate, func(tx context.Context, actor principal.Principal) error {
		id, err := validator.NormalizeID(roleID)
		if err != nil {
			return err
		}
		normalized, err := validator.NormalizeUpdateMetadata(request)
		if err != nil {
			return err
		}
		if err := service.inspector.Require(tx, actor.TenantID, id, normalized.Version); err != nil {
			return err
		}
		before, err := service.snapshot(tx, actor.TenantID, id)
		if err != nil {
			return err
		}
		result, err = service.Service.UpdateMetadata(tx, id, normalized)
		if err != nil {
			return err
		}
		after, err := service.snapshot(tx, actor.TenantID, id)
		if err != nil {
			return err
		}
		return service.recorder.Record(tx, &before, after)
	})
	if err != nil {
		return dto.Detail{}, err
	}
	return result, nil
}

func (service *AuditedService) snapshot(ctx context.Context, tenantID string, roleID string) (changed.Snapshot, error) {
	detail, err := service.repository.FindByID(ctx, tenantID, roleID)
	if err != nil {
		return changed.Snapshot{}, err
	}
	grants, err := service.repository.Permissions(ctx, tenantID, roleID)
	if err != nil {
		return changed.Snapshot{}, err
	}
	if detail.Role.TenantID != tenantID || detail.Role.RoleID != roleID || grants.RoleID != roleID || detail.Role.Version != grants.Version {
		return changed.Snapshot{}, repo.ErrVersionConflict
	}
	return changed.Snapshot{ID: roleID, Version: detail.Role.Version, Role: &changed.Role{Name: detail.Role.Name, Description: detail.Role.Description, Status: detail.Role.Status, Personas: append([]persona.Code{}, detail.Personas...), Permissions: append([]permission.Code{}, grants.Codes...)}}, nil
}
