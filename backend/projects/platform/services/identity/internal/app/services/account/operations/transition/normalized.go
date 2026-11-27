package transition

import (
	"context"
	"slices"
	"time"

	accountdto "grovehub/services/identity/internal/app/dto/account"
	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/entity"
	accountproducer "grovehub/services/identity/internal/app/events/producers/account"
	accountmapper "grovehub/services/identity/internal/app/mappers/account"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	"grovehub/services/identity/internal/app/repository/accountrole"
	"grovehub/services/identity/internal/app/repository/authority"
	"grovehub/services/identity/internal/app/repository/rbacmigration"
	sessionrepo "grovehub/services/identity/internal/app/repository/session"
	"grovehub/services/identity/internal/app/services/account/helpers"
	"grovehub/services/identity/internal/app/services/authorization"
	"grovehub/services/identity/internal/app/services/security/principal"
)

type Request struct {
	AccountID string
	Action string
	Version int64
}

type Normalized struct {
	guard authority.Guard
	readiness rbacmigration.Readiness
	authorizer authorization.Service
	accounts accountrepo.Reader
	writer accountrepo.StatusWriter
	assignments accountrole.Repository
	sessions sessionrepo.Repository
	mapper accountmapper.Mapper
	events accountproducer.ChangedProducer
}

// NewNormalized is intentionally not wired into live routes. The approval
// executor must bind this operation to the reviewed request and independent decision.
func NewNormalized(guard authority.Guard, readiness rbacmigration.Readiness, authorizer authorization.Service, accounts accountrepo.Reader, writer accountrepo.StatusWriter, assignments accountrole.Repository, sessions sessionrepo.Repository, mapper accountmapper.Mapper, events accountproducer.ChangedProducer) *Normalized {
	return &Normalized{guard: guard, readiness: readiness, authorizer: authorizer, accounts: accounts, writer: writer, assignments: assignments, sessions: sessions, mapper: mapper, events: events}
}

func (operation *Normalized) Execute(ctx context.Context, request Request) (accountdto.Account, error) {
	actor, err := principal.RequireTenant(ctx)
	if err != nil {
		return accountdto.Account{}, err
	}
	code, known := transitionPermission(request.Action)
	if !known || request.Version < 1 {
		return accountdto.Account{}, accountrepo.ErrConflict
	}
	if operation.guard == nil || operation.readiness == nil || operation.authorizer == nil || operation.accounts == nil || operation.writer == nil || operation.assignments == nil || operation.sessions == nil || operation.mapper == nil || operation.events == nil {
		return accountdto.Account{}, authorization.ErrUnavailable
	}
	var updated accountdto.Account
	err = operation.guard.Run(ctx, actor.TenantID, func(tx context.Context) error {
		if err := operation.readiness.Require(tx, actor.TenantID); err != nil {
			return err
		}
		effective, err := operation.authorizer.Effective(tx)
		if err != nil {
			return err
		}
		if !slices.Contains(effective, code) {
			return authorization.ErrForbidden
		}
		current, err := operation.accounts.FindByID(tx, request.AccountID)
		if err != nil {
			return err
		}
		if current.TenantID == nil || *current.TenantID != actor.TenantID {
			return accountrepo.ErrNotFound
		}
		if current.Version != request.Version {
			return accountrepo.ErrVersionConflict
		}
		next, err := helpers.NextStatus(current.Status, request.Action)
		if err != nil {
			return err
		}
		if err := operation.checkRoles(tx, actor.TenantID, current, next, effective); err != nil {
			return err
		}
		if err := operation.writer.ChangeStatus(tx, accountrepo.StatusChange{TenantID: actor.TenantID, AccountID: current.AccountID, Version: current.Version, Before: current.Status, After: next, UpdatedAt: time.Now().UTC()}); err != nil {
			return err
		}
		if next == entity.Suspended || next == entity.Deactivated {
			if _, err := operation.sessions.RevokeAll(tx, current.AccountID); err != nil {
				return err
			}
		}
		after, err := operation.accounts.FindByID(tx, current.AccountID)
		if err != nil {
			return err
		}
		if err := operation.events.Enqueue(tx, eventOperation(next), &current, &after); err != nil {
			return err
		}
		updated = operation.mapper.ToResponse(after)
		return nil
	})
	if err != nil {
		return accountdto.Account{}, err
	}
	return updated, nil
}

func (operation *Normalized) checkRoles(ctx context.Context, tenantID string, current entity.Account, next entity.AccountStatus, effective []permission.Code) error {
	assignment, err := operation.assignments.Find(ctx, tenantID, current.AccountID)
	if err != nil {
		return err
	}
	if assignment.Version != current.Version {
		return accountrepo.ErrVersionConflict
	}
	ids := slices.Clone(assignment.SecondaryRoleIDs)
	if assignment.PrimaryRoleID != "" {
		ids = append(ids, assignment.PrimaryRoleID)
	} else if next == entity.Active || len(ids) > 0 {
		return accountrole.ErrIneligible
	}
	candidates, err := operation.assignments.Candidates(ctx, tenantID, ids)
	if err != nil {
		return err
	}
	if len(candidates) != len(ids) {
		return accountrole.ErrIneligible
	}
	seen := make(map[string]bool, len(ids))
	for _, candidate := range candidates {
		if seen[candidate.RoleID] || !slices.Contains(ids, candidate.RoleID) {
			return accountrole.ErrIneligible
		}
		seen[candidate.RoleID] = true
		if next == entity.Active && (candidate.Status != entity.RoleActive || len(candidate.Codes) == 0) {
			return accountrole.ErrIneligible
		}
		for _, code := range candidate.Codes {
			if !permission.SchoolDelegable(code) || !slices.Contains(effective, code) {
				return authorization.ErrForbidden
			}
		}
	}
	return nil
}

func transitionPermission(action string) (permission.Code, bool) {
	switch action {
	case "suspend":
		return permission.AccountSuspend, true
	case "activate":
		return permission.AccountActivate, true
	case "deactivate":
		return permission.AccountDeactivate, true
	default:
		return "", false
	}
}

func eventOperation(status entity.AccountStatus) string {
	if status == entity.Active {
		return "ACTIVATED"
	}
	return string(status)
}
