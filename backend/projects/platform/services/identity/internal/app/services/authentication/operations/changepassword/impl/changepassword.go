package impl

import (
	"context"
	coredatabase "grovehub/libs/core/database"
	sessionrepo "grovehub/services/identity/internal/app/repository/session"
	changepassword "grovehub/services/identity/internal/app/services/authentication/operations/changepassword"
	credential "grovehub/services/identity/internal/app/services/credential"
	"grovehub/services/identity/internal/app/services/security/principal"
	"time"
)

type Service struct {
	sessions   sessionrepo.Repository
	credential credential.Service
	unit       coredatabase.UnitOfWork
}

var _ changepassword.Operation = (*Service)(nil)

func New(sessions sessionrepo.Repository, credentials credential.Service, unit coredatabase.UnitOfWork) changepassword.Operation {
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
