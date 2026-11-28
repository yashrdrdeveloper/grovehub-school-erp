package impl

import (
	"context"
	schoolgateway "grovehub/internal/session"
	"time"

	dto "grovehub/services/identity/internal/app/dto/account"
	"grovehub/services/identity/internal/app/entity"
	producer "grovehub/services/identity/internal/app/events/producers/account"
	repo "grovehub/services/identity/internal/app/repository/account"
	"grovehub/services/identity/internal/app/repository/accountrole"
	authrepo "grovehub/services/identity/internal/app/repository/authorization"
	"grovehub/services/identity/internal/app/repository/authority"
	tokens "grovehub/services/identity/internal/app/repository/passwordtoken"
	"grovehub/services/identity/internal/app/repository/rbacmigration"
	"grovehub/services/identity/internal/app/services/account/operations/setup"
	"grovehub/services/identity/internal/app/services/account/operations/setup/eligibility"
	secrets "grovehub/services/identity/internal/app/services/authentication/helpers"
	"grovehub/services/identity/internal/app/services/credential"
	"grovehub/services/identity/internal/app/services/authorization"
)

type Normalized struct {
	accounts repo.Reader
	writer repo.SetupWriter
	tokens tokens.Repository
	credentials credential.Service
	locker authority.Locker
	readiness rbacmigration.Readiness
	assignments accountrole.Repository
	snapshots authrepo.Repository
	events producer.ChangedProducer
}

func NewNormalized(accounts repo.Reader, writer repo.SetupWriter, tokens tokens.Repository, credentials credential.Service, locker authority.Locker, readiness rbacmigration.Readiness, assignments accountrole.Repository, snapshots authrepo.Repository, events producer.ChangedProducer) setup.Operation {
	return &Normalized{accounts: accounts, writer: writer, tokens: tokens, credentials: credentials, locker: locker, readiness: readiness, assignments: assignments, snapshots: snapshots, events: events}
}

func (operation *Normalized) Execute(ctx context.Context, grant string, input dto.Setup) (dto.SetupResult, error) {
	if operation.accounts == nil || operation.writer == nil || operation.tokens == nil || operation.credentials == nil || operation.locker == nil || operation.readiness == nil || operation.events == nil {
		return dto.SetupResult{}, authorization.ErrUnavailable
	}
	digest, err := secrets.Digest(grant, 32)
	if err != nil {
		return dto.SetupResult{}, tokens.ErrInvalidToken
	}
	proof, err := operation.tokens.Inspect(ctx, digest, entity.SetupGrant, time.Now().UTC())
	if err != nil {
		return dto.SetupResult{}, err
	}
	if proof.AccountVersion == nil {
		return dto.SetupResult{}, tokens.ErrInvalidToken
	}
	initial, err := operation.accounts.FindByID(ctx, proof.AccountID)
	if err != nil || initial.TenantID == nil || *initial.TenantID == "" {
		return dto.SetupResult{}, tokens.ErrInvalidToken
	}
	tenantID := *initial.TenantID
	if schoolgateway.CheckTenantSession(ctx, tenantID, proof.SessionGeneration) != nil { return dto.SetupResult{}, tokens.ErrInvalidToken }
	if err := operation.tokens.ReserveSetupAttempt(ctx, digest, time.Now().UTC()); err != nil {
		return dto.SetupResult{}, err
	}
	var result dto.SetupResult
	err = operation.locker.Run(ctx, tenantID, func(tx context.Context) error {
		if err := operation.readiness.Require(tx, tenantID); err != nil {
			return err
		}
		current, err := operation.accounts.FindByID(tx, proof.AccountID)
		if err != nil {
			return err
		}
		if current.TenantID == nil || *current.TenantID != tenantID || current.Version != *proof.AccountVersion {
			return tokens.ErrInvalidToken
		}
		if schoolgateway.CheckTenantSession(tx, tenantID, proof.SessionGeneration) != nil { return tokens.ErrInvalidToken }
		if err := eligibility.Require(tx, current, operation.assignments, operation.snapshots); err != nil {
			return err
		}
		id, err := operation.tokens.Consume(tx, digest, entity.SetupGrant, time.Now().UTC())
		if err != nil {
			return err
		}
		if id != current.AccountID {
			return tokens.ErrInvalidToken
		}
		if err := operation.credentials.Set(tx, id, input.Password); err != nil {
			return err
		}
		if err := operation.writer.ActivateSetup(tx, tenantID, id, current.Version, time.Now().UTC()); err != nil {
			return err
		}
		after, err := operation.accounts.FindByID(tx, id)
		if err != nil {
			return err
		}
		if err := operation.tokens.Invalidate(tx, id, entity.AccountSetup); err != nil {
			return err
		}
		if err := operation.tokens.Invalidate(tx, id, entity.SetupGrant); err != nil {
			return err
		}
		if err := operation.events.Enqueue(tx, "ACTIVATED", &current, &after); err != nil {
			return err
		}
		result = dto.SetupResult{AccountID: id, PersonGID: after.PersonGID, Status: string(after.Status), SetupCompleted: true}
		return nil
	})
	if err != nil {
		return dto.SetupResult{}, err
	}
	return result, nil
}
