package impl

import (
    "context"
    "encoding/json"
    "errors"

    "grovehub/services/integration/internal/app/entity"
    "grovehub/services/integration/internal/app/providers"
    connectionrepo "grovehub/services/integration/internal/app/repository/connection"
    connectionhelpers "grovehub/services/integration/internal/app/services/connection/helpers"
    actionstrategy "grovehub/services/integration/internal/app/services/connection/operations/action/strategies"
    oauth "grovehub/services/integration/internal/app/services/connection/operations/oauth"
)

type actionStrategy struct {
    repository connectionrepo.ConnectionRepository
    actions providers.ActionRegistry
    tokens oauth.TokenProvider
}

var _ actionstrategy.Strategy = (*actionStrategy)(nil)

func New(repository connectionrepo.ConnectionRepository, actions providers.ActionRegistry, tokens oauth.TokenProvider) actionstrategy.Strategy {
    return &actionStrategy{repository:repository,actions:actions,tokens:tokens}
}

func (s *actionStrategy) ExecuteAction(ctx context.Context, tenantID, connectionID string, action providers.Action, input json.RawMessage) (any, error) {
	c, err := s.repository.Get(ctx, tenantID, connectionID)
	if err != nil {
		return nil, err
	}
	if c.Status != entity.StatusActive {
		return nil, connectionhelpers.ErrConnectionInactive
	}
	if s.actions == nil {
		return nil, providers.ErrProviderUnavailable
	}
	if len(s.actions.Actions(c.Provider)) == 0 {
		return nil, providers.ErrProviderUnavailable
	}
	token, err := s.tokens.ValidProviderToken(ctx, connectionID, c.Provider)
	if err != nil {
		return nil, err
	}
	if token.AccessToken == "" {
		return nil, errors.New("provider authentication unavailable")
	}
	return s.actions.Execute(ctx, providers.Session{
		ConnectionID: connectionID, ProviderID: c.Provider, AccessToken: token.AccessToken,
	}, action, input)
}
