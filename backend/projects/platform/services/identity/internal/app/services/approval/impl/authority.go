package impl

import (
	"context"
	"slices"

	kind "grovehub/services/identity/internal/app/constants/approval"
	"grovehub/services/identity/internal/app/constants/permission"
	dto "grovehub/services/identity/internal/app/dto/approval"
	"grovehub/services/identity/internal/app/entity"
	repo "grovehub/services/identity/internal/app/repository/approval"
	accountrepo "grovehub/services/identity/internal/app/repository/accountrole"
	rolerepo "grovehub/services/identity/internal/app/repository/role"
	"grovehub/services/identity/internal/app/services/approval/helpers"
	"grovehub/services/identity/internal/app/services/accountrole/operations/preview"
	"grovehub/services/identity/internal/app/services/authorization"
	"grovehub/services/identity/internal/app/services/authorization/operations/delegation"
	"grovehub/services/identity/internal/app/services/security/principal"
)

func (service *Service) CreatePermissions(ctx context.Context, request dto.CreatePermissions) (entity.ApprovalRequest, error) {
	actor, err := principal.RequireTenant(ctx)
	if err != nil {
		return entity.ApprovalRequest{}, err
	}
	if !service.available() {
		return entity.ApprovalRequest{}, authorization.ErrUnavailable
	}
	intent, err := helpers.NormalizePermissions(request.Intent)
	if err != nil {
		return entity.ApprovalRequest{}, err
	}
	reason, err := helpers.NormalizeReason(request.Reason, request.IdempotencyKey)
	if err != nil {
		return entity.ApprovalRequest{}, err
	}
	payload, digest, err := helpers.EncodeIntent(intent)
	if err != nil {
		return entity.ApprovalRequest{}, err
	}
	return service.create(ctx, actor, kind.RolePermissionChange, permission.RolePermissionManage, reason, request.IdempotencyKey, payload, digest, func(tx context.Context) error {
		return service.validatePermissions(tx, actor.TenantID, intent)
	})
}

func (service *Service) validatePermissions(ctx context.Context, tenantID string, intent dto.PermissionIntent) error {
	detail, err := service.roles.FindByID(ctx, tenantID, intent.RoleID)
	if err != nil {
		return err
	}
	if detail.Role.Version != intent.Version {
		return rolerepo.ErrVersionConflict
	}
	// Only the existing draft permission writer can eventually apply this
	// intent. Live role changes need a separate guarded execution primitive.
	if detail.Role.Status != entity.RoleInactive {
		return rolerepo.ErrLiveRole
	}
	current, err := service.roles.Permissions(ctx, tenantID, intent.RoleID)
	if err != nil {
		return err
	}
	if current.Version != intent.Version {
		return rolerepo.ErrVersionConflict
	}
	effective, err := service.authorizer.Effective(ctx)
	if err != nil {
		return err
	}
	change, err := delegation.Preview(effective, current.Codes, intent.Codes)
	if err != nil {
		return err
	}
	if len(change.Added)+len(change.Removed) == 0 {
		return repo.ErrInvalid
	}
	return nil
}

func (service *Service) CreateAssignment(ctx context.Context, request dto.CreateAssignment) (entity.ApprovalRequest, error) {
	actor, err := principal.RequireTenant(ctx)
	if err != nil {
		return entity.ApprovalRequest{}, err
	}
	if !service.available() || service.assignments == nil {
		return entity.ApprovalRequest{}, authorization.ErrUnavailable
	}
	intent, err := helpers.NormalizeAssignment(request.Intent)
	if err != nil {
		return entity.ApprovalRequest{}, err
	}
	reason, err := helpers.NormalizeReason(request.Reason, request.IdempotencyKey)
	if err != nil {
		return entity.ApprovalRequest{}, err
	}
	payload, digest, err := helpers.EncodeIntent(intent)
	if err != nil {
		return entity.ApprovalRequest{}, err
	}
	return service.create(ctx, actor, kind.RoleAssignmentChange, permission.AccountRoleAssign, reason, request.IdempotencyKey, payload, digest, func(tx context.Context) error {
		return service.validateAssignment(tx, actor.TenantID, intent)
	})
}

func (service *Service) validateAssignment(ctx context.Context, tenantID string, intent dto.AssignmentIntent) error {
	current, err := service.assignments.Find(ctx, tenantID, intent.AccountID)
	if err != nil {
		return err
	}
	request := helpers.AssignmentReplacement(intent)
	candidates, err := service.assignments.Candidates(ctx, tenantID, preview.RoleIDs(current, request))
	if err != nil {
		return err
	}
	effective, err := service.authorizer.Effective(ctx)
	if err != nil {
		return err
	}
	result, err := preview.Evaluate(current, request, candidates, effective)
	if err != nil {
		return err
	}
	if !slices.Equal(result.RoleVersions, intent.RoleVersions) {
		return accountrepo.ErrVersionConflict
	}
	if !result.ApprovalRequired {
		return repo.ErrInvalid
	}
	return nil
}
