package impl

import (
	"context"
	"errors"
	coreevents "grovehub/libs/core/events"
	"grovehub/services/identity/internal/app/events/producers/resource"

	kind "grovehub/services/identity/internal/app/constants/approval"
	accountdto "grovehub/services/identity/internal/app/dto/account"
	"grovehub/services/identity/internal/app/entity"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	"grovehub/services/identity/internal/app/repository/accountrole"
	repo "grovehub/services/identity/internal/app/repository/approval"
	"grovehub/services/identity/internal/app/repository/authority"
	authrepo "grovehub/services/identity/internal/app/repository/authorization"
	"grovehub/services/identity/internal/app/repository/rbacmigration"
	rolerepo "grovehub/services/identity/internal/app/repository/role"
	"grovehub/services/identity/internal/app/services/account/operations/transition"
	"grovehub/services/identity/internal/app/services/approval/helpers"
	"grovehub/services/identity/internal/app/services/approval/operations/state"
	"grovehub/services/identity/internal/app/services/authorization"
	"grovehub/services/identity/internal/app/services/security/principal"
)

type AccountLifecycle interface {
	Execute(ctx context.Context, request transition.Request) (accountdto.Account, error)
}

type RoleLifecycle interface {
	Execute(ctx context.Context, request rolerepo.Transition) (rolerepo.Detail, error)
}

type AssignmentEffect interface {
	Execute(ctx context.Context, request accountrole.Replacement) (accountrole.Assignment, error)
}

type PendingInvitation interface {
	Queue(ctx context.Context, tenantID string, assignment accountrole.Assignment) error
}

type Effects struct {
	PendingInvitation PendingInvitation
	Audit             resource.Recorder
	AccountLifecycle  AccountLifecycle
	RoleLifecycle     RoleLifecycle
	Permissions       rolerepo.DraftPermissionWriter
	Assignment        AssignmentEffect
}

type Executor struct {
	service   *Service
	guard     authority.Guard
	readiness rbacmigration.Readiness
	effects   Effects
}

// NewExecutor constructs the transactional effect executor. HTTP callers must
// pass the explicit tenant cutover gate and use composite policy readiness.
func NewExecutor(repository repo.Repository, locker authority.Locker, guard authority.Guard, readiness rbacmigration.Readiness, authorizer authorization.Service, snapshots authrepo.Repository, accounts accountrepo.Reader, roles rolerepo.Repository, assignments accountrole.Repository, effects Effects) *Executor {
	return &Executor{service: &Service{repository: repository, locker: locker, authorizer: authorizer, snapshots: snapshots, accounts: accounts, roles: roles, assignments: assignments}, guard: guard, readiness: readiness, effects: effects}
}

func (executor *Executor) Execute(ctx context.Context, requestID string, version int64) (entity.ApprovalRequest, error) {
	actor, err := principal.RequireTenant(ctx)
	if err != nil {
		return entity.ApprovalRequest{}, err
	}
	if executor.service == nil || !executor.service.available() || executor.guard == nil || executor.readiness == nil || executor.effects.Audit == nil {
		return entity.ApprovalRequest{}, authorization.ErrUnavailable
	}
	id, err := helpers.NormalizeID(requestID)
	if err != nil || version < 1 {
		return entity.ApprovalRequest{}, repo.ErrInvalid
	}
	var updated entity.ApprovalRequest
	service := executor.service
	err = service.locker.Run(ctx, actor.TenantID, func(tx context.Context) error {
		current, err := service.repository.Lock(tx, actor.TenantID, id)
		if err != nil {
			return err
		}
		if current.RequesterID != actor.AccountID {
			return state.ErrForbidden
		}
		if err := helpers.VerifyIntent(current); err != nil {
			return err
		}
		// An already-applied request is a receipt, not a second mutation.
		if current.Status == kind.Applied {
			updated = current
			return nil
		}
		now, err := service.repository.Now(tx)
		if err != nil {
			return err
		}
		next, err := state.EffectResult(helpers.Snapshot(current), version, kind.Applied, now)
		if errors.Is(err, state.ErrExpired) {
			expired, err := state.Expire(helpers.Snapshot(current), version, now)
			if err != nil {
				return err
			}
			updated = helpers.Advance(current, expired, now)
			return service.repository.Advance(tx, current, updated)
		}
		if err != nil {
			return err
		}
		if err := executor.readiness.Require(tx, actor.TenantID); err != nil {
			return err
		}
		if err := executor.reviewers(tx, current); err != nil {
			return err
		}
		// Preflight must be read-only. Only deterministic version conflicts
		// become STALE; all errors after writing must abort the outer tx.
		apply, err := executor.prepare(tx, actor, current)
		if errors.Is(err, accountrepo.ErrVersionConflict) || errors.Is(err, rolerepo.ErrVersionConflict) || errors.Is(err, accountrole.ErrVersionConflict) {
			stale, err := state.EffectResult(helpers.Snapshot(current), version, kind.Stale, now)
			if err != nil {
				return err
			}
			updated = helpers.Advance(current, stale, now)
			return service.repository.Advance(tx, current, updated)
		}
		if err != nil {
			return err
		}
		return executor.guard.Run(tx, actor.TenantID, func(effectCtx context.Context) error {
			effectCtx = coreevents.WithCausationID(effectCtx, current.RequestID)
			before, err := executor.resourceSnapshot(effectCtx, current)
			if err != nil {
				return err
			}
			if err := apply(effectCtx); err != nil {
				return err
			}
			after, err := executor.resourceSnapshot(effectCtx, current)
			if err != nil {
				return err
			}
			finished, err := service.repository.Now(effectCtx)
			if err != nil {
				return err
			}
			// Never commit an effect after its approval deadline.
			next, err = state.EffectResult(helpers.Snapshot(current), version, kind.Applied, finished)
			if err != nil {
				return err
			}
			if err := executor.effects.Audit.Record(effectCtx, current, before, after, finished); err != nil {
				return err
			}
			updated = helpers.Advance(current, next, finished)
			return service.repository.Advance(effectCtx, current, updated)
		})
	})
	if err != nil {
		return entity.ApprovalRequest{}, err
	}
	return updated, nil
}

func (executor *Executor) reviewers(ctx context.Context, request entity.ApprovalRequest) error {
	records, err := executor.service.repository.Decisions(ctx, request.TenantID, request.RequestID)
	if err != nil {
		return err
	}
	if len(records) != request.StepCount {
		return repo.ErrConflict
	}
	for index, record := range records {
		if record.Step != index+1 || record.RequesterID != request.RequesterID || record.ReviewerID == request.RequesterID || record.Decision != kind.Approve || record.At.Before(request.CreatedAt) || !record.At.Before(request.ExpiresAt) {
			return repo.ErrConflict
		}
		if err := helpers.ReviewerIndependent(request, record.ReviewerID); err != nil {
			return err
		}
		step, err := executor.service.repository.Step(ctx, request.TenantID, request.RequestID, record.Step)
		if err != nil {
			return err
		}
		if err := executor.service.requireReviewer(ctx, principal.Principal{TenantID: request.TenantID, AccountID: record.ReviewerID}, step); err != nil {
			return err
		}
	}
	return nil
}
