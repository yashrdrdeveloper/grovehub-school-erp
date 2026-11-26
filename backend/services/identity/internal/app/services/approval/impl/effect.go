package impl

import (
	"context"
	"encoding/json"
	"slices"

	kind "grovehub/services/identity/internal/app/constants/approval"
	"grovehub/services/identity/internal/app/constants/permission"
	dto "grovehub/services/identity/internal/app/dto/approval"
	"grovehub/services/identity/internal/app/entity"
	"grovehub/services/identity/internal/app/repository/accountrole"
	repo "grovehub/services/identity/internal/app/repository/approval"
	rolerepo "grovehub/services/identity/internal/app/repository/role"
	"grovehub/services/identity/internal/app/services/account/operations/transition"
	"grovehub/services/identity/internal/app/services/approval/helpers"
	"grovehub/services/identity/internal/app/services/authorization"
	"grovehub/services/identity/internal/app/services/security/principal"
)

func (executor *Executor) prepare(ctx context.Context, actor principal.Principal, request entity.ApprovalRequest) (func(context.Context) error, error) {
	switch request.RequestType {
	case kind.AccountLifecycleChange, kind.RoleLifecycleChange:
		return executor.prepareLifecycle(ctx, actor, request)
	case kind.RolePermissionChange:
		return executor.preparePermissions(ctx, actor, request)
	case kind.RoleAssignmentChange:
		return executor.prepareAssignment(ctx, actor, request)
	default:
		return nil, repo.ErrInvalid
	}
}

func (executor *Executor) prepareLifecycle(ctx context.Context, actor principal.Principal, request entity.ApprovalRequest) (func(context.Context) error, error) {
	var intent dto.LifecycleIntent
	if err := json.Unmarshal(request.Payload, &intent); err != nil {
		return nil, err
	}
	normalized, code, _, _, err := helpers.NormalizeLifecycle(dto.CreateLifecycle{Type: request.RequestType, Intent: intent, Reason: request.Reason, IdempotencyKey: request.IdempotencyKey})
	if err != nil {
		return nil, err
	}
	if err := executor.service.authorizer.Require(ctx, code); err != nil {
		return nil, err
	}
	if err := executor.service.validateTarget(ctx, actor.TenantID, normalized); err != nil {
		return nil, err
	}
	intent = normalized.Intent
	if request.RequestType == kind.AccountLifecycleChange {
		if executor.effects.AccountLifecycle == nil {
			return nil, authorization.ErrUnavailable
		}
		return func(tx context.Context) error {
			_, err := executor.effects.AccountLifecycle.Execute(tx, transition.Request{AccountID: intent.TargetID, Action: intent.Action, Version: intent.Version})
			return err
		}, nil
	}
	if executor.effects.RoleLifecycle == nil {
		return nil, authorization.ErrUnavailable
	}
	grants, err := executor.service.roles.Permissions(ctx, actor.TenantID, intent.TargetID)
	if err != nil {
		return nil, err
	}
	if grants.Version != intent.Version {
		return nil, rolerepo.ErrVersionConflict
	}
	effective, err := executor.service.authorizer.Effective(ctx)
	if err != nil {
		return nil, err
	}
	for _, grant := range grants.Codes {
		if !permission.SchoolDelegable(grant) || !slices.Contains(effective, grant) {
			return nil, authorization.ErrForbidden
		}
	}
	status := entity.RoleInactive
	if intent.Action == "activate" {
		status = entity.RoleActive
	}
	return func(tx context.Context) error {
		_, err := executor.effects.RoleLifecycle.Execute(tx, rolerepo.Transition{TenantID: actor.TenantID, ActorID: actor.AccountID, RoleID: intent.TargetID, Version: intent.Version, Status: status})
		return err
	}, nil
}

func (executor *Executor) preparePermissions(ctx context.Context, actor principal.Principal, request entity.ApprovalRequest) (func(context.Context) error, error) {
	if executor.effects.Permissions == nil {
		return nil, authorization.ErrUnavailable
	}
	var intent dto.PermissionIntent
	if err := json.Unmarshal(request.Payload, &intent); err != nil {
		return nil, err
	}
	intent, err := helpers.NormalizePermissions(intent)
	if err != nil {
		return nil, err
	}
	if err := executor.service.validatePermissions(ctx, actor.TenantID, intent); err != nil {
		return nil, err
	}
	return func(tx context.Context) error {
		_, err := executor.effects.Permissions.ReplaceDraftPermissions(tx, rolerepo.DraftPermissionUpdate{TenantID: actor.TenantID, ActorID: actor.AccountID, RoleID: intent.RoleID, Version: intent.Version, Codes: intent.Codes})
		return err
	}, nil
}

func (executor *Executor) prepareAssignment(ctx context.Context, actor principal.Principal, request entity.ApprovalRequest) (func(context.Context) error, error) {
	if executor.effects.Assignment == nil || executor.service.assignments == nil {
		return nil, authorization.ErrUnavailable
	}
	var intent dto.AssignmentIntent
	if err := json.Unmarshal(request.Payload, &intent); err != nil {
		return nil, err
	}
	intent, err := helpers.NormalizeAssignment(intent)
	if err != nil {
		return nil, err
	}
	if err := executor.service.validateAssignment(ctx, actor.TenantID, intent); err != nil {
		return nil, err
	}
	versions := make([]accountrole.RoleVersion, 0, len(intent.RoleVersions))
	for _, item := range intent.RoleVersions {
		versions = append(versions, accountrole.RoleVersion{RoleID: item.RoleID, Version: item.Version})
	}
	return func(tx context.Context) error {
		assignment, err := executor.effects.Assignment.Execute(tx, accountrole.Replacement{TenantID: actor.TenantID, ActorID: actor.AccountID, AccountID: intent.AccountID, Version: intent.Version, PrimaryRoleID: intent.PrimaryRoleID, SecondaryRoleIDs: intent.SecondaryRoleIDs, RoleVersions: versions})
		if err != nil {
			return err
		}
		if executor.effects.PendingInvitation != nil {
			return executor.effects.PendingInvitation.Queue(tx, actor.TenantID, assignment)
		}
		return nil
	}, nil
}
