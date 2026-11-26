package impl

import (
	"context"

	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/constants/persona"
	"grovehub/services/identity/internal/app/entity"
	"grovehub/services/identity/internal/app/events/producers/resource/operations/changed"
	repo "grovehub/services/identity/internal/app/repository/approval"
	"grovehub/services/identity/internal/app/services/approval/helpers"
	"grovehub/services/identity/internal/app/services/authorization"
)

func (executor *Executor) resourceSnapshot(ctx context.Context, request entity.ApprovalRequest) (changed.Snapshot, error) {
	id, role, err := helpers.ResourceID(request)
	if err != nil {
		return changed.Snapshot{}, err
	}
	if role {
		detail, err := executor.service.roles.FindByID(ctx, request.TenantID, id)
		if err != nil {
			return changed.Snapshot{}, err
		}
		grants, err := executor.service.roles.Permissions(ctx, request.TenantID, id)
		if err != nil {
			return changed.Snapshot{}, err
		}
		if detail.Role.RoleID != id || detail.Role.TenantID != request.TenantID || grants.RoleID != id || grants.Version != detail.Role.Version {
			return changed.Snapshot{}, repo.ErrConflict
		}
		return changed.Snapshot{ID: id, Version: detail.Role.Version, Role: &changed.Role{Name: detail.Role.Name, Description: detail.Role.Description, Status: detail.Role.Status, Personas: append([]persona.Code{}, detail.Personas...), Permissions: append([]permission.Code{}, grants.Codes...)}}, nil
	}
	if executor.service.assignments == nil {
		return changed.Snapshot{}, authorization.ErrUnavailable
	}
	account, err := executor.service.assignments.Find(ctx, request.TenantID, id)
	if err != nil {
		return changed.Snapshot{}, err
	}
	if account.AccountID != id {
		return changed.Snapshot{}, repo.ErrConflict
	}
	return changed.Snapshot{ID: id, Version: account.Version, Account: &changed.Account{Status: account.Status, PrimaryRoleID: account.PrimaryRoleID, SecondaryRoleIDs: append([]string{}, account.SecondaryRoleIDs...)}}, nil
}
