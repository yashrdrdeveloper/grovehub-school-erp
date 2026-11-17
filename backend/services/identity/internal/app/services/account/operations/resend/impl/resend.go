package impl

import (
	"context"
	accountdto "grovehub/services/identity/internal/app/dto/account"
	"errors"
	"grovehub/services/identity/internal/app/entity"
	"grovehub/services/identity/internal/app/security/principal"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	registration "grovehub/services/identity/internal/app/services/registration"
	transaction "grovehub/libs/core/database/transaction"
	resend "grovehub/services/identity/internal/app/services/account/operations/resend"
)

type Service struct {
	accounts accountrepo.Repository
	invitations registration.Service
	unit transaction.UnitOfWork
}

var _ resend.Operation = (*Service)(nil)

func New(accounts accountrepo.Repository, invitations registration.Service, unit transaction.UnitOfWork) resend.Operation {
	return &Service{accounts: accounts, invitations: invitations, unit: unit}
}

func (op *Service) Execute(ctx context.Context, accountID string) error {
	actor, err := principal.RequireAdmin(ctx)
	if err != nil {
		return err
	}

	return op.unit.Run(ctx, func(transactionCtx context.Context) error {
		accountEntity, err := op.accounts.FindByID(transactionCtx, accountID)
		if err != nil {
			return err
		}
		if accountEntity.TenantID == nil || *accountEntity.TenantID != actor.TenantID {
			return accountrepo.ErrNotFound
		}
		if accountEntity.Status != entity.PendingSetup {
			return errors.New("account is not pending setup")
		}
		return op.invitations.Queue(transactionCtx, accountEntity, entity.AccountSetup, "there")
	})
}
