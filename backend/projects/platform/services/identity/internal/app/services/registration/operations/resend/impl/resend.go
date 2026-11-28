package impl

import (
	"context"
	"errors"
	coredatabase "grovehub/libs/core/database"
	"grovehub/services/identity/internal/app/entity"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	registration "grovehub/services/identity/internal/app/services/registration"
	resend "grovehub/services/identity/internal/app/services/registration/operations/resend"
	"grovehub/services/identity/internal/app/services/security/principal"
)

type Service struct {
	accounts    accountrepo.Repository
	invitations registration.InvitationSender
	unit        coredatabase.UnitOfWork
}

var _ resend.Operation = (*Service)(nil)

func New(accounts accountrepo.Repository, invitations registration.InvitationSender, unit coredatabase.UnitOfWork) resend.Operation {
	return &Service{accounts: accounts, invitations: invitations, unit: unit}
}

func (op *Service) Execute(ctx context.Context, accountID string) error {
	actor, err := principal.RequireTenant(ctx)
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
