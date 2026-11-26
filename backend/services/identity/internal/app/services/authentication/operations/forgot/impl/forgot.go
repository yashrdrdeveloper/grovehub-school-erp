package impl

import (
	"context"
	"errors"
	"strings"
	"grovehub/services/identity/internal/app/entity"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	registration "grovehub/services/identity/internal/app/services/registration"
	transaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
	forgot "grovehub/services/identity/internal/app/services/authentication/operations/forgot"
)

type Service struct {
	accounts accountrepo.Repository
	invitation registration.InvitationSender
	unit transaction.UnitOfWork
}

var _ forgot.Operation = (*Service)(nil)

func New(accounts accountrepo.Repository, invitations registration.InvitationSender, unit transaction.UnitOfWork) forgot.Operation {
	return &Service{accounts: accounts, invitation: invitations, unit: unit}
}

func (op *Service) Execute(ctx context.Context, email string) error {
	email = strings.ToLower(strings.TrimSpace(email))
	if !strings.Contains(email, "@") {
		return nil
	}

	account, err := op.accounts.FindByEmail(ctx, email)
	if err != nil || account.Status != entity.Active {
		return nil
	}
	if op.invitation == nil {
		return errors.New("mail unavailable")
	}

	return op.unit.Run(ctx, func(transactionCtx context.Context) error {
		return op.invitation.Queue(transactionCtx, account, entity.PasswordReset, "there")
	})
}
