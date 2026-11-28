package impl

import (
	"context"
	"time"
	"grovehub/services/identity/internal/app/services/security/principal"
	sessionrepo "grovehub/services/identity/internal/app/repository/session"
	logoutall "grovehub/services/identity/internal/app/services/authentication/operations/logoutall"
)

type Service struct {
	sessions sessionrepo.Repository
}

var _ logoutall.Operation = (*Service)(nil)

func New(sessions sessionrepo.Repository) logoutall.Operation {
	return &Service{sessions: sessions}
}

func (op *Service) Execute(ctx context.Context) (int64, error) {
	actor, err := principal.Require(ctx)
	if err != nil {
		return 0, err
	}
	if err := principal.Recent(ctx, 5*time.Minute); err != nil {
		return 0, err
	}
	return op.sessions.RevokeAll(ctx, actor.AccountID)
}
