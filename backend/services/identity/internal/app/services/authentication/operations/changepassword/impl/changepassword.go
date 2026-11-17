package impl

import (
	"context"
	"time"
	"grovehub/services/identity/internal/app/security/principal"
	sessionrepo "grovehub/services/identity/internal/app/repository/session"
	credential "grovehub/services/identity/internal/app/services/credential"
	transaction "grovehub/libs/core/database/transaction"
	changepassword "grovehub/services/identity/internal/app/services/authentication/operations/changepassword"
)

type Service struct {
	sessions sessionrepo.Repository
	credential credential.Service
	unit transaction.UnitOfWork
}

var _ changepassword.Operation = (*Service)(nil)

func New(sessions sessionrepo.Repository, credentials credential.Service, unit transaction.UnitOfWork) changepassword.Operation {
	return &Service{sessions: sessions, credential: credentials, unit: unit}
}

func (op *Service) Execute(ctx context.Context, password string) error {
	actor, err := principal.Require(ctx)
	if err != nil {
		return err
	}
	if err := principal.Recent(ctx, 5*time.Minute); err != nil {
		return err
	}

	return op.unit.Run(ctx, func(transactionCtx context.Context) error {
		if err := op.credential.Replace(transactionCtx, actor.AccountID, password); err != nil {
			return err
		}
		_, err := op.sessions.RevokeAll(transactionCtx, actor.AccountID)
		return err
	})
}
