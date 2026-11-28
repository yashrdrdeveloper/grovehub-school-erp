package impl

import (
	"context"

	"grovehub/services/identity/internal/app/entity"
	"grovehub/services/identity/internal/app/services/approval/helpers"
	"grovehub/services/identity/internal/app/services/approval/operations/state"
	"grovehub/services/identity/internal/app/services/authorization"
	"grovehub/services/identity/internal/app/services/security/principal"
)

func (service *Service) Cancel(ctx context.Context, requestID string, version int64) (entity.ApprovalRequest, error) {
	actor, err := principal.RequireTenant(ctx)
	if err != nil {
		return entity.ApprovalRequest{}, err
	}
	if !service.available() {
		return entity.ApprovalRequest{}, authorization.ErrUnavailable
	}
	id, err := helpers.NormalizeID(requestID)
	if err != nil {
		return entity.ApprovalRequest{}, err
	}
	var updated entity.ApprovalRequest
	err = service.locker.Run(ctx, actor.TenantID, func(tx context.Context) error {
		current, err := service.repository.Lock(tx, actor.TenantID, id)
		if err != nil {
			return err
		}
		if current.RequesterID != actor.AccountID {
			return state.ErrForbidden
		}
		snapshot, err := service.snapshots.LoadSnapshot(tx, actor.TenantID, actor.AccountID)
		if err != nil {
			return err
		}
		if snapshot.Status != entity.Active || snapshot.AccountID != actor.AccountID || snapshot.TenantID != actor.TenantID {
			return authorization.ErrForbidden
		}
		now, err := service.repository.Now(tx)
		if err != nil {
			return err
		}
		next, err := state.Cancel(helpers.Snapshot(current), state.Actor{TenantID: actor.TenantID, AccountID: actor.AccountID}, version, now)
		if err != nil {
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
