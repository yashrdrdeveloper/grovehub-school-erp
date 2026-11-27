package impl

import (
    providerdto "grovehub/services/integration/internal/app/providers/dto"
    providercontracts "grovehub/services/integration/internal/app/providers/contracts"
    providerhelpers "grovehub/services/integration/internal/app/providers/helpers"
    "context"
    "encoding/json"
    "errors"
    "github.com/google/uuid"
    "grovehub/services/integration/internal/app/services/connection/operations/action/delivery"

    "grovehub/services/integration/internal/app/entity"
    mailconstants "grovehub/services/integration/internal/app/providers/mailprovider/constants"
    connectionrepo "grovehub/services/integration/internal/app/repository/connection"
    connectionhelpers "grovehub/services/integration/internal/app/services/connection/helpers"
    actionstrategy "grovehub/services/integration/internal/app/services/connection/operations/action/strategies"
    oauth "grovehub/services/integration/internal/app/services/connection/operations/oauth"
)

type actionStrategy struct {
    repository connectionrepo.ConnectionRepository
    actions providercontracts.ActionRegistry
    tokens oauth.TokenProvider
    delivery *delivery.Dispatcher
}

var _ actionstrategy.Strategy = (*actionStrategy)(nil)

func New(repository connectionrepo.ConnectionRepository, actions providercontracts.ActionRegistry, tokens oauth.TokenProvider, dispatchers ...*delivery.Dispatcher) actionstrategy.Strategy {
    strategy := &actionStrategy{repository:repository,actions:actions,tokens:tokens}
    if len(dispatchers) > 0 { strategy.delivery = dispatchers[0] }
    return strategy
}

func (s *actionStrategy) ExecuteAction(ctx context.Context, tenantID, connectionID string, action providerdto.Action, input json.RawMessage) (any, error) {
	return s.executeAction(ctx, tenantID, connectionID, action, input, "")
}

func (s *actionStrategy) executeAction(ctx context.Context, tenantID, connectionID string, action providerdto.Action, input json.RawMessage, deliveryID string) (any, error) {
	c, err := s.repository.Get(ctx, tenantID, connectionID)
	if err != nil {
		return nil, err
	}
	if c.Status != entity.StatusActive {
		return nil, connectionhelpers.ErrConnectionInactive
	}
	if s.actions == nil {
		return nil, providerhelpers.ErrProviderUnavailable
	}
	if len(s.actions.Actions(c.Provider)) == 0 {
		return nil, providerhelpers.ErrProviderUnavailable
	}
	token, err := s.tokens.ValidProviderToken(ctx, connectionID, c.Provider)
	if err != nil {
		return nil, err
	}
	if token.AccessToken == "" {
		return nil, errors.New("provider authentication unavailable")
	}
	send := func() (any, error) {
		return s.actions.Execute(ctx, providerdto.Session{
			ConnectionID: connectionID, ProviderID: c.Provider, AccessToken: token.AccessToken,
		}, action, input)
	}
	if deliveryID != "" {
		return s.delivery.Dispatch(ctx, tenantID, deliveryID, input, send)
	}
	return send()
}

func (s *actionStrategy) ExecuteSystemMail(ctx context.Context, tenantID string, input json.RawMessage) (any, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(input, &fields); err != nil {
		return nil, err
	}
	var deliveryID string
	if raw, ok := fields["deliveryId"]; ok {
		if err := json.Unmarshal(raw, &deliveryID); err != nil {
			return nil, errors.New("invalid mail delivery identifier")
		}
		parsed, err := uuid.Parse(deliveryID)
		if err != nil || parsed == uuid.Nil || parsed.String() != deliveryID || s.delivery == nil {
			return nil, errors.New("mail delivery deduplication unavailable or invalid identifier")
		}
		delete(fields, "deliveryId")
		input, err = json.Marshal(fields)
		if err != nil { return nil, err }
		found, err := s.delivery.Replay(ctx, tenantID, deliveryID, input)
		if err != nil || found { return nil, err }
	}
	sender, err := s.repository.SystemMailSender(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	return s.executeAction(ctx, tenantID, sender.ID, mailconstants.MailSend, input, deliveryID)
}

func (s *actionStrategy) ConfigureSystemMailSender(ctx context.Context, tenantID string, connectionID string) error {
	sender, err := s.repository.Get(ctx, tenantID, connectionID)
	if err != nil {
		return err
	}
	if sender.Status != entity.StatusActive {
		return connectionhelpers.ErrConnectionInactive
	}
	if s.actions == nil {
		return providerhelpers.ErrProviderUnavailable
	}
	supported := false
	for _, action := range s.actions.Actions(sender.Provider) {
		if action == mailconstants.MailSend {
			supported = true
			break
		}
	}
	if !supported {
		return providerhelpers.ErrProviderUnavailable
	}
	return s.repository.ConfigureSystemMailSender(ctx, tenantID, connectionID)
}
