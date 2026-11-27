package impl

import (
	"context"
	"slices"

	"grovehub/services/identity/internal/app/constants/permission"
	roledto "grovehub/services/identity/internal/app/dto/role"
	rolerepo "grovehub/services/identity/internal/app/repository/role"
	"grovehub/services/identity/internal/app/services/authorization"
	"grovehub/services/identity/internal/app/services/authorization/operations/delegation"
	"grovehub/services/identity/internal/app/services/security/principal"
	rolevalidator "grovehub/services/identity/internal/app/validators/role"
)

func (service *Service) Permissions(ctx context.Context, roleID string) (roledto.PermissionSet, error) {
	actor, err := service.authorize(ctx, permission.RolePermissionView)
	if err != nil {
		return roledto.PermissionSet{}, err
	}
	roleID, err = rolevalidator.NormalizeID(roleID)
	if err != nil {
		return roledto.PermissionSet{}, err
	}
	set, err := service.repository.Permissions(ctx, actor.TenantID, roleID)
	if err != nil {
		return roledto.PermissionSet{}, err
	}
	return roledto.PermissionSet{RoleID: set.RoleID, Version: set.Version, Codes: append([]permission.Code{}, set.Codes...)}, nil
}

func (service *Service) PreviewPermissions(ctx context.Context, roleID string, request roledto.PreviewPermissions) (roledto.PermissionPreview, error) {
	actor, err := principal.RequireTenant(ctx)
	if err != nil {
		return roledto.PermissionPreview{}, err
	}
	if service.repository == nil || service.authorizer == nil {
		return roledto.PermissionPreview{}, authorization.ErrUnavailable
	}
	effective, err := service.authorizer.Effective(ctx)
	if err != nil {
		return roledto.PermissionPreview{}, err
	}
	if !slices.Contains(effective, permission.RolePermissionManage) {
		return roledto.PermissionPreview{}, authorization.ErrForbidden
	}
	roleID, err = rolevalidator.NormalizeID(roleID)
	if err != nil {
		return roledto.PermissionPreview{}, err
	}
	if request.Version < 1 {
		return roledto.PermissionPreview{}, rolevalidator.ErrInvalidVersion
	}
	if request.Codes == nil {
		return roledto.PermissionPreview{}, delegation.ErrInvalidSet
	}
	set, err := service.repository.Permissions(ctx, actor.TenantID, roleID)
	if err != nil {
		return roledto.PermissionPreview{}, err
	}
	if set.Version != request.Version {
		return roledto.PermissionPreview{}, rolerepo.ErrVersionConflict
	}
	change, err := delegation.Preview(effective, set.Codes, request.Codes)
	if err != nil {
		return roledto.PermissionPreview{}, err
	}
	return roledto.PermissionPreview{
		RoleID: set.RoleID, Version: set.Version, Added: change.Added, Removed: change.Removed,
		ApprovalRequired: len(change.Added)+len(change.Removed) > 0,
	}, nil
}
