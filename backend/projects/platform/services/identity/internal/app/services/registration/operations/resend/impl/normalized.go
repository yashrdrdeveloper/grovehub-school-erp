package impl

import (
	"context"
	"slices"
	"time"

	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/entity"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	"grovehub/services/identity/internal/app/repository/accountrole"
	"grovehub/services/identity/internal/app/repository/authority"
	"grovehub/services/identity/internal/app/repository/invitation"
	"grovehub/services/identity/internal/app/repository/rbacmigration"
	"grovehub/services/identity/internal/app/services/authorization"
	"grovehub/services/identity/internal/app/services/registration/operations/resend"
	"grovehub/services/identity/internal/app/services/security/principal"
)

type ApprovedSender interface {
	Queue(ctx context.Context, tenantID string, assignment accountrole.Assignment) error
}

type Normalized struct {
	locker authority.Locker
	readiness rbacmigration.Readiness
	authorizer authorization.Service
	assignments accountrole.Repository
	sender ApprovedSender
	cooldown invitation.CooldownReader
}

func NewNormalized(locker authority.Locker, readiness rbacmigration.Readiness, authorizer authorization.Service, assignments accountrole.Repository, sender ApprovedSender, cooldown invitation.CooldownReader) resend.Operation {
	return &Normalized{locker: locker, readiness: readiness, authorizer: authorizer, assignments: assignments, sender: sender, cooldown: cooldown}
}

func (operation *Normalized) Execute(ctx context.Context, accountID string) error {
	actor, err := principal.RequireTenant(ctx)
	if err != nil {
		return err
	}
	if operation.locker == nil || operation.readiness == nil || operation.authorizer == nil || operation.assignments == nil || operation.sender == nil || operation.cooldown == nil {
		return authorization.ErrUnavailable
	}
	return operation.locker.Run(ctx, actor.TenantID, func(tx context.Context) error {
		effective, err := operation.authorizer.Effective(tx)
		if err != nil {
			return err
		}
		if !slices.Contains(effective, permission.AccountResendInvitation) {
			return authorization.ErrForbidden
		}
		if err := operation.readiness.Require(tx, actor.TenantID); err != nil {
			return err
		}
		assignment, err := operation.assignments.Find(tx, actor.TenantID, accountID)
		if err != nil {
			return err
		}
		if assignment.Status != entity.PendingSetup || assignment.PrimaryRoleID == "" {
			return accountrepo.ErrConflict
		}
		ids := append([]string{assignment.PrimaryRoleID}, assignment.SecondaryRoleIDs...)
		candidates, err := operation.assignments.Candidates(tx, actor.TenantID, ids)
		if err != nil {
			return err
		}
		for _, candidate := range candidates {
			for _, code := range candidate.Codes {
				if !permission.SchoolDelegable(code) || !slices.Contains(effective, code) {
					return authorization.ErrForbidden
				}
			}
		}
		recent, err := operation.cooldown.RecentSetup(tx, assignment.AccountID, time.Now().UTC().Add(-30*time.Second))
		if err != nil {
			return err
		}
		if recent {
			return accountrepo.ErrConflict
		}
		return operation.sender.Queue(tx, actor.TenantID, assignment)
	})
}
