package impl

import (
	"context"
	"errors"
	coredatabase "grovehub/libs/core/database"
	schoolgateway "grovehub/libs/platform/session"
	accountdto "grovehub/services/identity/internal/app/dto/account"
	"grovehub/services/identity/internal/app/entity"
	accountproducer "grovehub/services/identity/internal/app/events/producers/account"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	tokenrepo "grovehub/services/identity/internal/app/repository/passwordtoken"
	accounthelpers "grovehub/services/identity/internal/app/services/account/helpers"
	setup "grovehub/services/identity/internal/app/services/account/operations/setup"
	secrets "grovehub/services/identity/internal/app/services/authentication/helpers"
	credential "grovehub/services/identity/internal/app/services/credential"
	"time"
)

type Service struct {
	accounts    accountrepo.Repository
	tokens      tokenrepo.Repository
	credentials credential.Service
	unit        coredatabase.UnitOfWork
	events      accountproducer.ChangedProducer
}

var _ setup.Operation = (*Service)(nil)

func New(accounts accountrepo.Repository, tokens tokenrepo.Repository, credentials credential.Service, unit coredatabase.UnitOfWork, events accountproducer.ChangedProducer) setup.Operation {
	return &Service{accounts: accounts, tokens: tokens, credentials: credentials, unit: unit, events: events}
}

func (op *Service) Execute(ctx context.Context, grant string, input accountdto.Setup) (accountdto.SetupResult, error) {
	digest, err := secrets.Digest(grant, 32)
	if err != nil {
		return accountdto.SetupResult{}, err
	}

	proof, err := op.tokens.Inspect(ctx, digest, entity.SetupGrant, time.Now().UTC())
	if err != nil {
		return accountdto.SetupResult{}, err
	}
	initial, err := op.accounts.FindByID(ctx, proof.AccountID)
	if err != nil || initial.Status != entity.PendingSetup || schoolgateway.CheckAccountGeneration(ctx, initial.TenantID, proof.SessionGeneration) != nil {
		return accountdto.SetupResult{}, tokenrepo.ErrInvalidToken
	}

	// The attempt is committed independently so transaction rollback cannot refund it.
	if err := op.tokens.ReserveSetupAttempt(ctx, digest, time.Now().UTC()); err != nil {
		return accountdto.SetupResult{}, err
	}

	var result accountdto.SetupResult
	err = op.unit.Run(ctx, func(transactionCtx context.Context) error {
		accountID, err := op.tokens.Consume(transactionCtx, digest, entity.SetupGrant, time.Now().UTC())
		if err != nil {
			return err
		}
		accountEntity, err := op.accounts.FindByID(transactionCtx, accountID)
		if err != nil {
			return err
		}
		if accountEntity.Status != entity.PendingSetup {
			return errors.New("account not pending setup")
		}

		if schoolgateway.CheckAccountGeneration(transactionCtx, accountEntity.TenantID, proof.SessionGeneration) != nil {
			return tokenrepo.ErrInvalidToken
		}
		before := accountEntity
		accountEntity.Status = entity.Active
		accountEntity.Version++
		accountEntity.UpdatedAt = time.Now().UTC()
		accountEntity.SearchDetails = accounthelpers.SearchDetails(accountEntity)

		if err := op.credentials.Set(transactionCtx, accountID, input.Password); err != nil {
			return err
		}
		if err := op.accounts.Update(transactionCtx, accountEntity); err != nil {
			return err
		}
		if err := op.events.Enqueue(transactionCtx, "ACTIVATED", &before, &accountEntity); err != nil {
			return err
		}

		result = accountdto.SetupResult{
			AccountID:      accountID,
			PersonGID:      accountEntity.PersonGID,
			Status:         string(accountEntity.Status),
			SetupCompleted: true,
		}
		return nil
	})
	return result, err
}
