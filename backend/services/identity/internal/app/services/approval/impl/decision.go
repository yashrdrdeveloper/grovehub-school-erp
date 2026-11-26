package impl

import (
	"context"

	dto "grovehub/services/identity/internal/app/dto/approval"
	"grovehub/services/identity/internal/app/entity"
	repo "grovehub/services/identity/internal/app/repository/approval"
	"grovehub/services/identity/internal/app/services/approval/helpers"
	"grovehub/services/identity/internal/app/services/approval/operations/state"
	"grovehub/services/identity/internal/app/services/authorization"
	"grovehub/services/identity/internal/app/services/authorization/operations/evaluate"
	"grovehub/services/identity/internal/app/services/security/principal"
)

func (service *Service) Decide(ctx context.Context, request dto.Review) (entity.ApprovalRequest, error) {
	actor, err := principal.RequireTenant(ctx)
	if err != nil {
		return entity.ApprovalRequest{}, err
	}
	if !service.available() {
		return entity.ApprovalRequest{}, authorization.ErrUnavailable
	}
	id, err := helpers.NormalizeID(request.RequestID)
	if err != nil {
		return entity.ApprovalRequest{}, err
	}
	var updated entity.ApprovalRequest
	err = service.locker.Run(ctx, actor.TenantID, func(tx context.Context) error {
		current, err := service.repository.Lock(tx, actor.TenantID, id)
		if err != nil {
			return err
		}
		if err := helpers.ReviewerIndependent(current, actor.AccountID); err != nil {
			return err
		}
		step, err := service.repository.Step(tx, actor.TenantID, current.RequestID, current.CurrentStep)
		if err != nil {
			return err
		}
		if err := service.requireReviewer(tx, actor, step); err != nil {
			return err
		}
		if err := helpers.VerifyIntent(current); err != nil {
			return err
		}
		now, err := service.repository.Now(tx)
		if err != nil {
			return err
		}
		next, err := state.Decide(helpers.Snapshot(current), state.Actor{TenantID: actor.TenantID, AccountID: actor.AccountID}, request.Version, request.Decision, request.Comment, now)
		if err != nil {
			return err
		}
		record := repo.ReviewRecord{RequesterID: current.RequesterID, ReviewerID: actor.AccountID, Decision: request.Decision, Comment: request.Comment, Step: current.CurrentStep, At: now}
		if err := service.repository.InsertDecision(tx, actor.TenantID, current.RequestID, record); err != nil {
			return err
		}
		updated = helpers.Advance(current, next, now)
		return service.repository.Advance(tx, current, updated)
	})
	if err != nil {
		return entity.ApprovalRequest{}, err
	}
	return updated, nil
}

func (service *Service) requireReviewer(ctx context.Context, actor principal.Principal, step repo.Step) error {
	snapshot, err := service.snapshots.LoadSnapshot(ctx, actor.TenantID, actor.AccountID)
	if err != nil {
		return err
	}
	if step.RoleID != nil {
		// Require the named active role itself to carry the step permission.
		roles := snapshot.Roles
		snapshot.Roles = nil
		for _, role := range roles {
			if role.RoleID == *step.RoleID {
				snapshot.Roles = append(snapshot.Roles, role)
			}
		}
	}
	if !evaluate.AllowedFor(actor.AccountID, actor.TenantID, step.Permission, snapshot) {
		return authorization.ErrForbidden
	}
	return nil
}
