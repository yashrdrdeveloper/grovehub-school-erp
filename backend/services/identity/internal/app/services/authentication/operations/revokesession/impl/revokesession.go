package impl

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"grovehub/services/identity/internal/app/services/security/principal"
	sessionrepo "grovehub/services/identity/internal/app/repository/session"
	revokesession "grovehub/services/identity/internal/app/services/authentication/operations/revokesession"
)

type Service struct {
	sessions sessionrepo.Repository
}

var _ revokesession.Operation = (*Service)(nil)

func New(sessions sessionrepo.Repository) revokesession.Operation {
	return &Service{sessions: sessions}
}

func (op *Service) Execute(ctx context.Context, sessionID string) error {
	actor, err := principal.Require(ctx)
	if err != nil {
		return err
	}
	if _, err := uuid.Parse(sessionID); err != nil {
		return errors.New("invalid session ID")
	}
	return op.sessions.Revoke(ctx, actor.AccountID, sessionID)
}
