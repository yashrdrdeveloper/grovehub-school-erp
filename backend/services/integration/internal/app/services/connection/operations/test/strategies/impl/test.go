package impl

import (
    providerdto "grovehub/services/integration/internal/app/providers/dto"
    providercontracts "grovehub/services/integration/internal/app/providers/contracts"
    providerhelpers "grovehub/services/integration/internal/app/providers/helpers"
    oauth "grovehub/services/integration/internal/app/providers/mailprovider/operations/oauth"
    mapper "grovehub/services/integration/internal/app/mappers/connection"
    "context"
    "encoding/json"
    "errors"
    "strings"
    "time"

    "github.com/google/uuid"
    connectionconstants "grovehub/services/integration/internal/app/constants/connections"
    connectiondto "grovehub/services/integration/internal/app/dto/connection"
    "grovehub/services/integration/internal/app/entity"
    maildto "grovehub/services/integration/internal/app/providers/mailprovider/dto"
    connectionrepo "grovehub/services/integration/internal/app/repository/connection"
    connectionhelpers "grovehub/services/integration/internal/app/services/connection/helpers"
    oauth "grovehub/services/integration/internal/app/services/connection/operations/oauth"
    teststrategies "grovehub/services/integration/internal/app/services/connection/operations/test/strategies"
)

type testStrategy struct {
    mapper mapper.ConnectionMapper
    repository connectionrepo.ConnectionRepository
    actions providercontracts.ActionRegistry
    tokens oauth.TokenProvider
}

var _ teststrategies.Strategy = (*testStrategy)(nil)

func New(repository connectionrepo.ConnectionRepository, actions providercontracts.ActionRegistry, tokens oauth.TokenProvider, eventMapper mapper.ConnectionMapper) teststrategies.Strategy {
    return &testStrategy{repository:repository,actions:actions,tokens:tokens,mapper:eventMapper}
}

func (s *testStrategy) RunTest(ctx context.Context, tenantID, id, direction string) (connectiondto.ConnectionTest, error) {
	connection, err := s.repository.Get(ctx, tenantID, id)
	if err != nil {
		return connectiondto.ConnectionTest{}, err
	}

	direction = strings.ToUpper(strings.TrimSpace(direction))
	if connection.Status != entity.StatusActive || connection.ExternalID == nil {
		return connectiondto.ConnectionTest{}, connectionhelpers.ErrConnectionInactive
	}
	action, err := connectionhelpers.ProbeAction(direction)
	if err != nil {
		return connectiondto.ConnectionTest{}, connectionhelpers.ErrInvalidConnection
	}

	started := time.Now().UTC()
	test := entity.ConnectionTest{
		ID:           uuid.NewString(),
		ConnectionID: id,
		Direction:    direction,
		Status:       connectionconstants.TestSuccess,
		StartedAt:    started,
	}
	payload := json.RawMessage("{}")

	if direction == connectionconstants.DirectionOutbound {
  payload, err = json.Marshal(maildto.Send{
   To: []string{*connection.ExternalID},
   Subject: "GrooveHub Connection Test",
   Body: maildto.Body{Text: "This email verifies that the connected provider accepts outgoing mail. Test ID: "+test.ID},
  })
 }
	if err == nil {
		var token maildto.OAuthToken
		token, err = s.tokens.ValidProviderToken(ctx, id, connection.Provider)
		if err == nil {
			if s.actions == nil {
				err = providerhelpers.ErrProviderUnavailable
			} else if token.AccessToken == "" {
				err = errors.New("provider token missing")
			} else {
				_, err = s.actions.Execute(ctx, providerdto.Session{
					ConnectionID: id,
					ProviderID:   connection.Provider,
					AccessToken:  token.AccessToken,
				}, action, payload)
			}
		}
	}
	if err != nil {
		test.Status = connectionconstants.TestFailed
		test.Detail = "Provider connection test failed; upstream details redacted"
	} else if direction == connectionconstants.DirectionInbound {
		test.Detail = "Provider mailbox read capability verified; inbound delivery not verified"
	} else {
		test.Detail = "Provider accepted self-addressed test message; delivery not confirmed"
	}

	test.CompletedAt = time.Now().UTC()
	if err := s.repository.SaveTest(ctx, test); err != nil {
		return connectiondto.ConnectionTest{}, err
	}
	return s.mapper.TestToResponse(test), nil
}

func (s *testStrategy) ListTests(ctx context.Context, tenantID, id string) ([]connectiondto.ConnectionTest, error) {
	if _, err := s.repository.Get(ctx, tenantID, id); err != nil {
		return nil, err
	}
	rows, err := s.repository.ListTests(ctx, id)
	if err != nil {
		return nil, err
	}
	result := make([]connectiondto.ConnectionTest, 0, len(rows))
	for _, test := range rows {
		result = append(result, s.mapper.TestToResponse(test))
	}
	return result, nil
}
