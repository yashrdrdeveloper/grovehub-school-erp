package impl

import (
	"context"

	kind "grovehub/services/identity/internal/app/constants/approval"
	dto "grovehub/services/identity/internal/app/dto/approval"
	"grovehub/services/identity/internal/app/entity"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	rolerepo "grovehub/services/identity/internal/app/repository/role"
	repo "grovehub/services/identity/internal/app/repository/approval"
	accounthelpers "grovehub/services/identity/internal/app/services/account/helpers"
)

func (service *Service) validateTarget(ctx context.Context, tenantID string, request dto.CreateLifecycle) error {
	switch request.Type {
	case kind.AccountLifecycleChange:
		account, err := service.accounts.FindByID(ctx, request.Intent.TargetID)
		if err != nil {
			return err
		}
		if account.TenantID == nil || *account.TenantID != tenantID {
			return accountrepo.ErrNotFound
		}
		if account.Version != request.Intent.Version {
			return accountrepo.ErrVersionConflict
		}
		_, err = accounthelpers.NextStatus(account.Status, request.Intent.Action)
		return err
	case kind.RoleLifecycleChange:
		role, err := service.roles.FindByID(ctx, tenantID, request.Intent.TargetID)
		if err != nil {
			return err
		}
		if role.Role.Version != request.Intent.Version {
			return rolerepo.ErrVersionConflict
		}
		if (request.Intent.Action == "activate" && role.Role.Status != entity.RoleInactive) || (request.Intent.Action == "deactivate" && role.Role.Status != entity.RoleActive) {
			return rolerepo.ErrInvalidTransition
		}
		return nil
	default:
		return repo.ErrInvalid
	}
}
