package impl

import (
	"context"
	accountdto "grovehub/services/identity/internal/app/dto/account"
	"errors"
	"time"
	"grovehub/services/identity/internal/app/entity"
	"grovehub/services/identity/internal/app/services/security/principal"
	accounthelpers "grovehub/services/identity/internal/app/services/account/helpers"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	accountmapper "grovehub/services/identity/internal/app/mappers/account"
	accountvalidator "grovehub/services/identity/internal/app/validators/account"
	transaction "grovehub/libs/core/database/transaction"
	accountproducer "grovehub/services/identity/internal/app/events/producers/account"
	update "grovehub/services/identity/internal/app/services/account/operations/update"
)

type Service struct {
	accounts accountrepo.Repository
	mapper accountmapper.Mapper
	validator accountvalidator.Validator
	unit transaction.UnitOfWork
	events accountproducer.ChangedProducer
}

var _ update.Operation = (*Service)(nil)

func New(accounts accountrepo.Repository, mapper accountmapper.Mapper, validator accountvalidator.Validator, unit transaction.UnitOfWork, events accountproducer.ChangedProducer) update.Operation {
	return &Service{accounts: accounts, mapper: mapper, validator: validator, unit: unit, events: events}
}

func (op *Service) Execute(ctx context.Context, accountID string, request accountdto.Update) (accountdto.Account, error) {
	actor, err := principal.RequireAdmin(ctx)
	if err != nil {
		return accountdto.Account{}, err
	}
	if request.Username == nil {
		return accountdto.Account{}, errors.New("username must be supplied")
	}
	username, err := op.validator.NormalizeUsername(*request.Username)
	if err != nil {
		return accountdto.Account{}, err
	}

	var updated entity.Account
	err = op.unit.Run(ctx, func(transactionCtx context.Context) error {
		current, err := op.accounts.FindByID(transactionCtx, accountID)
		if err != nil {
			return err
		}
		if current.TenantID == nil || *current.TenantID != actor.TenantID {
			return accountrepo.ErrNotFound
		}

		before := current
		current.Username = &username
		current.Version++
		current.UpdatedAt = time.Now().UTC()
		current.SearchDetails = accounthelpers.SearchDetails(current)
		if err := op.accounts.Update(transactionCtx, current); err != nil {
			return err
		}
		if err := op.events.Enqueue(transactionCtx, "UPDATED", &before, &current); err != nil {
			return err
		}
		updated = current
		return nil
	})
	if err != nil {
		return accountdto.Account{}, err
	}
	return op.mapper.ToResponse(updated), nil
}
