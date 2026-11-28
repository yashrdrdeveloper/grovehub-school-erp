package impl

import (
	"context"
	"grovehub/services/identity/internal/app/services/security/principal"
	authdto "grovehub/services/identity/internal/app/dto/authentication"
	sessionrepo "grovehub/services/identity/internal/app/repository/session"
	sessions "grovehub/services/identity/internal/app/services/authentication/operations/sessions"
)

type Service struct {
	sessions sessionrepo.Repository
}

var _ sessions.Operation = (*Service)(nil)

func New(sessions sessionrepo.Repository) sessions.Operation {
	return &Service{sessions: sessions}
}

func (op *Service) Execute(ctx context.Context) (authdto.SessionsResult, error) {
	actor, err := principal.Require(ctx)
	if err != nil {
		return authdto.SessionsResult{}, err
	}
	sessions, err := op.sessions.List(ctx, actor.AccountID)
	if err != nil {
		return authdto.SessionsResult{}, err
	}

	items := make([]authdto.SessionSummary, 0, len(sessions))
	for _, session := range sessions {
		items = append(items, authdto.SessionSummary{
			SessionID: session.SessionID,
			Current: session.SessionID == actor.SessionID,
			CreatedAt: session.CreatedAt,
			LastActiveAt: session.LastActiveAt,
			ExpiresAt: session.RefreshExpiresAt,
		})
	}
	return authdto.SessionsResult{Items: items}, nil
}
