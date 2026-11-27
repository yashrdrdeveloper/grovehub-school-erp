package impl

import (
	"context"
	"math"
	"slices"
	"time"

	"github.com/google/uuid"

	"grovehub/services/identity/internal/app/constants/permission"
	dto "grovehub/services/identity/internal/app/dto/account"
	"grovehub/services/identity/internal/app/entity"
	"grovehub/services/identity/internal/app/events/producers/accountmetadata"
	mapper "grovehub/services/identity/internal/app/mappers/account"
	repo "grovehub/services/identity/internal/app/repository/account"
	"grovehub/services/identity/internal/app/repository/accountrole"
	"grovehub/services/identity/internal/app/repository/authority"
	tokens "grovehub/services/identity/internal/app/repository/passwordtoken"
	"grovehub/services/identity/internal/app/repository/rbacmigration"
	"grovehub/services/identity/internal/app/services/account/helpers"
	"grovehub/services/identity/internal/app/services/account/operations/update"
	"grovehub/services/identity/internal/app/services/authorization"
	"grovehub/services/identity/internal/app/services/security/principal"
	validator "grovehub/services/identity/internal/app/validators/account"
)

type Normalized struct {
	accounts repo.Reader
	writer repo.MetadataWriter
	assignments accountrole.Repository
	tokens tokens.Repository
	locker authority.Locker
	readiness rbacmigration.Readiness
	authorizer authorization.Service
	mapper mapper.Mapper
	validator validator.Validator
	audit accountmetadata.Recorder
}

func NewNormalized(accounts repo.Reader, writer repo.MetadataWriter, assignments accountrole.Repository, tokens tokens.Repository, locker authority.Locker, readiness rbacmigration.Readiness, authorizer authorization.Service, mapper mapper.Mapper, validator validator.Validator, audit accountmetadata.Recorder) update.Operation {
	return &Normalized{accounts: accounts, writer: writer, assignments: assignments, tokens: tokens, locker: locker, readiness: readiness, authorizer: authorizer, mapper: mapper, validator: validator, audit: audit}
}

func (operation *Normalized) Execute(ctx context.Context, accountID string, request dto.Update) (dto.Account, error) {
	actor, err := principal.RequireTenant(ctx)
	if err != nil {
		return dto.Account{}, err
	}
	if operation.accounts == nil || operation.writer == nil || operation.assignments == nil || operation.tokens == nil || operation.locker == nil || operation.readiness == nil || operation.authorizer == nil || operation.mapper == nil || operation.validator == nil || operation.audit == nil {
		return dto.Account{}, authorization.ErrUnavailable
	}
	id, err := uuid.Parse(accountID)
	if err != nil || id == uuid.Nil {
		return dto.Account{}, repo.ErrNotFound
	}
	if request.RoleAssignment != nil || request.Username == nil || request.Version == nil || *request.Version < 1 || *request.Version == math.MaxInt64 {
		return dto.Account{}, repo.ErrConflict
	}
	username, err := operation.validator.NormalizeUsername(*request.Username)
	if err != nil {
		return dto.Account{}, err
	}
	var result dto.Account
	err = operation.locker.Run(ctx, actor.TenantID, func(tx context.Context) error {
		effective, err := operation.authorizer.Effective(tx)
		if err != nil {
			return err
		}
		if !slices.Contains(effective, permission.AccountUpdate) {
			return authorization.ErrForbidden
		}
		if err := operation.readiness.Require(tx, actor.TenantID); err != nil {
			return err
		}
		current, err := operation.accounts.FindByID(tx, id.String())
		if err != nil {
			return err
		}
		if current.TenantID == nil || *current.TenantID != actor.TenantID {
			return repo.ErrNotFound
		}
		if err := operation.requireDelegation(tx, actor.TenantID, current, effective); err != nil {
			return err
		}
		if current.Version != *request.Version {
			return repo.ErrVersionConflict
		}
		if current.Status != entity.Active && current.Status != entity.PendingSetup {
			return repo.ErrConflict
		}
		if current.Username != nil && *current.Username == username {
			result = operation.mapper.ToResponse(current)
			return nil
		}
		after := current
		after.Username = &username
		after.SearchDetails = helpers.SearchDetails(after)
		if err := operation.writer.UpdateMetadata(tx, repo.MetadataUpdate{TenantID: actor.TenantID, AccountID: current.AccountID, Version: current.Version, Status: current.Status, Username: username, SearchDetails: after.SearchDetails, UpdatedAt: time.Now().UTC()}); err != nil {
			return err
		}
		if current.Status == entity.PendingSetup {
			for _, purpose := range []entity.TokenPurpose{entity.AccountSetup, entity.SetupGrant} {
				if err := operation.tokens.Invalidate(tx, current.AccountID, purpose); err != nil {
					return err
				}
			}
		}
		after, err = operation.accounts.FindByID(tx, current.AccountID)
		if err != nil {
			return err
		}
		if err := operation.audit.Record(tx, current, after); err != nil {
			return err
		}
		result = operation.mapper.ToResponse(after)
		return nil
	})
	if err != nil {
		return dto.Account{}, err
	}
	return result, nil
}

func (operation *Normalized) requireDelegation(ctx context.Context, tenantID string, current entity.Account, effective []permission.Code) error {
	assignment, err := operation.assignments.Find(ctx, tenantID, current.AccountID)
	if err != nil {
		return err
	}
	if assignment.Version != current.Version || assignment.Status != current.Status {
		return repo.ErrVersionConflict
	}
	ids := slices.Clone(assignment.SecondaryRoleIDs)
	if assignment.PrimaryRoleID != "" {
		ids = append(ids, assignment.PrimaryRoleID)
	} else if current.Status != entity.PendingSetup || len(ids) != 0 {
		return accountrole.ErrIneligible
	}
	candidates, err := operation.assignments.Candidates(ctx, tenantID, ids)
	if err != nil {
		return err
	}
	if len(candidates) != len(ids) {
		return accountrole.ErrIneligible
	}
	for _, candidate := range candidates {
		for _, code := range candidate.Codes {
			if !permission.SchoolDelegable(code) || !slices.Contains(effective, code) {
				return authorization.ErrForbidden
			}
		}
	}
	return nil
}
