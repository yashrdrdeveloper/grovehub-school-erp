package impl

import (
	"context"
	"grovehub/services/identity/internal/app/services/security/principal"
	sessionrepo "grovehub/services/identity/internal/app/repository/session"
	logout "grovehub/services/identity/internal/app/services/authentication/operations/logout"
)

type Service struct {
	sessions sessionrepo.Repository
}

var _ logout.Operation = (*Service)(nil)

func New(sessions sessionrepo.Repository) logout.Operation {
	return &Service{sessions: sessions}
}

func (op *Service) Execute(ctx context.Context) error {
	actor, err := principal.Require(ctx)
	if err != nil {
		return err
	}
	return op.sessions.Revoke(ctx, actor.AccountID, actor.SessionID)
}
