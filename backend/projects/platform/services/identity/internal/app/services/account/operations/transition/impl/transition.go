package impl

import (
	"context"
	coredatabase "grovehub/libs/core/database"
	accountdto "grovehub/services/identity/internal/app/dto/account"
	"grovehub/services/identity/internal/app/entity"
	accountproducer "grovehub/services/identity/internal/app/events/producers/account"
	accountmapper "grovehub/services/identity/internal/app/mappers/account"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	sessionrepo "grovehub/services/identity/internal/app/repository/session"
	"grovehub/services/identity/internal/app/services/account/helpers"
	transition "grovehub/services/identity/internal/app/services/account/operations/transition"
	"grovehub/services/identity/internal/app/services/security/principal"
	"time"
)

type Service struct {
	accounts accountrepo.Repository
	sessions sessionrepo.Repository
	mapper   accountmapper.Mapper
	unit     coredatabase.UnitOfWork
	events   accountproducer.ChangedProducer
}

var _ transition.Operation = (*Service)(nil)

func New(accounts accountrepo.Repository, sessions sessionrepo.Repository, mapper accountmapper.Mapper, unit coredatabase.UnitOfWork, events accountproducer.ChangedProducer) transition.Operation {
	return &Service{accounts: accounts, sessions: sessions, mapper: mapper, unit: unit, events: events}
}

func (op *Service) Execute(ctx context.Context, accountID string, action string) (accountdto.Account, error) {
	actor, err := principal.RequireTenant(ctx)
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

		next, err := helpers.NextStatus(current.Status, action)
		if err != nil {
			return err
		}

		before := current
		current.Status = next
		current.Version++
		current.UpdatedAt = time.Now().UTC()
		if err := op.accounts.Update(transactionCtx, current); err != nil {
			return err
		}

		if next == entity.Suspended || next == entity.Deactivated {
			if _, err := op.sessions.RevokeAll(transactionCtx, current.AccountID); err != nil {
				return err
			}
		}

		operation := map[string]string{
			"suspend":    "SUSPENDED",
			"activate":   "ACTIVATED",
			"deactivate": "DEACTIVATED",
		}[action]
		if err := op.events.Enqueue(transactionCtx, operation, &before, &current); err != nil {
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
