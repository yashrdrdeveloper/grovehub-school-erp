package impl

import (
	"context"
	"errors"
	"time"
	"grovehub/services/identity/internal/app/services/security/principal"
	sessionrepo "grovehub/services/identity/internal/app/repository/session"
	credential "grovehub/services/identity/internal/app/services/credential"
	reauthenticate "grovehub/services/identity/internal/app/services/authentication/operations/reauthenticate"
)

type Service struct {
	sessions sessionrepo.Repository
	credential credential.Service
}

var _ reauthenticate.Operation = (*Service)(nil)

func New(sessions sessionrepo.Repository, credentials credential.Service) reauthenticate.Operation {
	return &Service{sessions: sessions, credential: credentials}
}

func (op *Service) Execute(ctx context.Context, password string) error {
	actor, err := principal.Require(ctx)
	if err != nil {
		return err
	}
	valid, err := op.credential.Verify(ctx, actor.AccountID, password)
	if err != nil || !valid {
		return errors.New("invalid credentials")
	}
	return op.sessions.Reauthenticate(ctx, actor.SessionID, time.Now().UTC())
}
