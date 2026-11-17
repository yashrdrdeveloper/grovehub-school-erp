package impl

import (
	"context"
	accountdto "grovehub/services/identity/internal/app/dto/account"
	"errors"
	"time"
	"grovehub/services/identity/internal/app/entity"
	secrets "grovehub/services/identity/internal/app/services/authentication/helpers"
	accounthelpers "grovehub/services/identity/internal/app/services/account/helpers"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	tokenrepo "grovehub/services/identity/internal/app/repository/passwordtoken"
	credential "grovehub/services/identity/internal/app/services/credential"
	accountvalidator "grovehub/services/identity/internal/app/validators/account"
	transaction "grovehub/libs/core/database/transaction"
	accountproducer "grovehub/services/identity/internal/app/events/producers/account"
	setup "grovehub/services/identity/internal/app/services/account/operations/setup"
)

type Service struct {
	accounts accountrepo.Repository
	tokens tokenrepo.Repository
	credentials credential.Service
	validator accountvalidator.Validator
	unit transaction.UnitOfWork
	events accountproducer.ChangedProducer
}

var _ setup.Operation = (*Service)(nil)

func New(accounts accountrepo.Repository, tokens tokenrepo.Repository, credentials credential.Service, validator accountvalidator.Validator, unit transaction.UnitOfWork, events accountproducer.ChangedProducer) setup.Operation {
	return &Service{accounts: accounts, tokens: tokens, credentials: credentials, validator: validator, unit: unit, events: events}
}

func (op *Service) Execute(ctx context.Context, grant string, input accountdto.Setup) (accountdto.SetupResult, error) {
	username, err := op.validator.NormalizeUsername(input.Username)
	if err != nil {
		return accountdto.SetupResult{}, err
	}
	digest, err := secrets.Digest(grant, 32)
	if err != nil {
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

		before := accountEntity
		accountEntity.Username = &username
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
			AccountID: accountID,
			PersonGID: accountEntity.PersonGID,
			Username: username,
			Status: string(accountEntity.Status),
			SetupCompleted: true,
		}
		return nil
	})
	return result, err
}
