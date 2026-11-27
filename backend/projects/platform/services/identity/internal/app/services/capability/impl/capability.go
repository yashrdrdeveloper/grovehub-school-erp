package impl

import (
	"context"
	"slices"

	"grovehub/services/identity/internal/app/constants/permission"
	dto "grovehub/services/identity/internal/app/dto/capability"
	"grovehub/services/identity/internal/app/repository/authority"
	repo "grovehub/services/identity/internal/app/repository/capability"
	"grovehub/services/identity/internal/app/services/approval/operations/cutover"
	"grovehub/services/identity/internal/app/services/authorization"
	"grovehub/services/identity/internal/app/services/capability"
	"grovehub/services/identity/internal/app/services/capability/helpers"
	"grovehub/services/identity/internal/app/services/security/principal"
)

type Service struct {
	reader repo.Reader
	locker authority.Locker
	authorizer authorization.Service
	executionTenant string
}

var _ capability.Service = (*Service)(nil)

func New(reader repo.Reader, locker authority.Locker, authorizer authorization.Service, executionTenant string) capability.Service {
	return &Service{reader: reader, locker: locker, authorizer: authorizer, executionTenant: executionTenant}
}

func (service *Service) Discover(ctx context.Context) (dto.Discovery, error) {
	actor, err := principal.RequireTenant(ctx)
	if err != nil {
		return dto.Discovery{}, err
	}
	if service.locker == nil || service.authorizer == nil {
		return dto.Discovery{}, authorization.ErrUnavailable
	}
	var result dto.Discovery
	err = service.locker.Run(ctx, actor.TenantID, func(tx context.Context) error {
		effective, err := service.authorizer.Effective(tx)
		if err != nil {
			return err
		}
		result.Permissions = append([]permission.Code{}, effective...)
		result.Actions = helpers.Actions(effective)
		result.ExecutionConfigured = service.executionTenant != "" && service.executionTenant == actor.TenantID
		result.ExecutionRequiresReauthentication = result.ExecutionConfigured && cutover.New(service.executionTenant).Require(tx) != nil
		result.Actions["accountCreate"] = result.ExecutionConfigured && slices.Contains(effective, permission.AccountCreate)
		result.Actions["accountUpdate"] = result.ExecutionConfigured && slices.Contains(effective, permission.AccountUpdate)
		result.Actions["accountResendInvitation"] = result.ExecutionConfigured && slices.Contains(effective, permission.AccountResendInvitation)
		return nil
	})
	if err != nil {
		return dto.Discovery{}, err
	}
	return result, nil
}

func (service *Service) Catalog(ctx context.Context) (dto.Catalog, error) {
	actor, err := principal.RequireTenant(ctx)
	if err != nil {
		return dto.Catalog{}, err
	}
	if service.reader == nil || service.locker == nil || service.authorizer == nil {
		return dto.Catalog{}, authorization.ErrUnavailable
	}
	var result dto.Catalog
	err = service.locker.Run(ctx, actor.TenantID, func(tx context.Context) error {
		effective, err := service.authorizer.Effective(tx)
		if err != nil {
			return err
		}
		if !slices.Contains(effective, permission.PermissionView) {
			return authorization.ErrForbidden
		}
		rows, err := service.reader.Available(tx, actor.TenantID)
		if err != nil {
			return err
		}
		result = helpers.Catalog(rows, effective)
		return nil
	})
	if err != nil {
		return dto.Catalog{}, err
	}
	return result, nil
}
