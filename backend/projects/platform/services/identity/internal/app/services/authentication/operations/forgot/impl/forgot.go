package impl

import (
	"context"
	schoolgateway "grovehub/libs/platform/session"
	"errors"
	"strings"
	"grovehub/services/identity/internal/app/entity"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	registration "grovehub/services/identity/internal/app/services/registration"
	coredatabase "grovehub/libs/core/database"
	forgot "grovehub/services/identity/internal/app/services/authentication/operations/forgot"
)

type Service struct {
	accounts accountrepo.Repository
	invitation registration.InvitationSender
	unit coredatabase.UnitOfWork
}

var _ forgot.Operation = (*Service)(nil)

func New(accounts accountrepo.Repository, invitations registration.InvitationSender, unit coredatabase.UnitOfWork) forgot.Operation {
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
		current, err := op.accounts.FindByID(transactionCtx, account.AccountID)
		if err != nil || current.Status != entity.Active || current.Email != email || schoolgateway.CheckAccountTenant(transactionCtx, current.TenantID) != nil { return nil }
		return op.invitation.Queue(transactionCtx, current, entity.PasswordReset, "there")
	})
}
