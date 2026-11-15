package impl

import (
    "context"
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "encoding/hex"
    "errors"
    "time"

    connectiondto "grovehub/services/integration/internal/app/dto/connection"
    "grovehub/services/integration/internal/app/entity"
    "grovehub/services/integration/internal/app/providers"
    connectionrepo "grovehub/services/integration/internal/app/repository/connection"
    connectionhelpers "grovehub/services/integration/internal/app/services/connection/helpers"
    strategies "grovehub/services/integration/internal/app/services/connection/operations/oauth/strategies"
)

type oauthStrategy struct {
    repository connectionrepo.ConnectionRepository
    oauth providers.OAuthRegistry
    cipher connectionhelpers.TokenCipher
}
var _ strategies.Strategy = (*oauthStrategy)(nil)
func New(repository connectionrepo.ConnectionRepository, oauth providers.OAuthRegistry, cipher connectionhelpers.TokenCipher) strategies.Strategy {
    return &oauthStrategy{repository:repository,oauth:oauth,cipher:cipher}
}

func (s *oauthStrategy) BeginOAuth(ctx context.Context, tenantID, id string) (connectiondto.OAuthStart, error) {
	c, err := s.repository.Get(ctx, tenantID, id)
	if err != nil {
		return connectiondto.OAuthStart{}, err
	}
	if c.Status != entity.StatusPendingSetup {
		return connectiondto.OAuthStart{}, connectionhelpers.ErrInvalidConnection
	}
	strategy, err := s.oauth.Find(c.Provider)
	if err != nil {
		return connectiondto.OAuthStart{}, connectionhelpers.ErrOAuthUnavailable
	}
	nonce := make([]byte, 32)
	if _, err = rand.Read(nonce); err != nil {
		return connectiondto.OAuthStart{}, err
	}
	state := base64.RawURLEncoding.EncodeToString(nonce)
	hash := hashState(state)
	err = s.repository.SaveOAuthState(ctx, entity.OAuthState{
		Hash: hash, ConnectionID: id, ExpiresAt: time.Now().UTC().Add(10 * time.Minute),
	})
	if err != nil {
		return connectiondto.OAuthStart{}, err
	}
	return connectiondto.OAuthStart{AuthorizationURL: strategy.AuthorizationURL(state)}, nil
}

func (s *oauthStrategy) FinishOAuth(ctx context.Context, provider, code, state string) (connectiondto.OAuthCompletion, error) {
	if code == "" || len(state) < 40 {
		return connectiondto.OAuthCompletion{}, connectionhelpers.ErrOAuthInvalid
	}
	oauthState, err := s.repository.ConsumeOAuthState(ctx, hashState(state), time.Now().UTC())
	if err != nil {
		if errors.Is(err, integrationrepository.ErrOAuthStateInvalid) {
			return connectiondto.OAuthCompletion{}, connectionhelpers.ErrOAuthInvalid
		}
		return connectiondto.OAuthCompletion{}, err
	}
	connection, err := s.repository.GetByID(ctx, oauthState.ConnectionID)
	if err != nil {
		return connectiondto.OAuthCompletion{}, err
	}
	if connection.Provider != provider || connection.Status != entity.StatusPendingSetup {
		return connectiondto.OAuthCompletion{}, connectionhelpers.ErrOAuthInvalid
	}
	strategy, err := s.oauth.Find(provider)
	if err != nil {
		return connectiondto.OAuthCompletion{}, connectionhelpers.ErrOAuthUnavailable
	}
	token, err := strategy.Exchange(ctx, code)
	if err != nil {
		return connectiondto.OAuthCompletion{}, err
	}
	if token.RefreshToken == "" {
		return connectiondto.OAuthCompletion{}, connectionhelpers.ErrOAuthInvalid
	}
	mailbox, err := strategy.Mailbox(ctx, token)
	if err != nil {
		return connectiondto.OAuthCompletion{}, err
	}
	if mailbox == "" || token.RefreshToken == "" {
		return connectiondto.OAuthCompletion{}, connectionhelpers.ErrOAuthInvalid
	}
	encrypted, err := s.cipher.Seal(token)
	if err != nil {
		return connectiondto.OAuthCompletion{}, err
	}
	if err = s.repository.Activate(ctx, connection.ID, provider, mailbox, encrypted); err != nil {
		return connectiondto.OAuthCompletion{}, err
	}
	return connectiondto.OAuthCompletion{Status: entity.StatusActive}, nil
}

func hashState(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func (s *oauthStrategy) ValidProviderToken(ctx context.Context, connectionID, provider string) (providers.OAuthToken, error) {
	var empty providers.OAuthToken
	encrypted, err := s.repository.Credential(ctx, connectionID)
	if err != nil {
		return empty, err
	}
	token, err := s.cipher.Open(encrypted)
	if err != nil {
		return empty, err
	}
	if time.Until(token.Expiry) > time.Minute {
		return token, nil
	}
	strategy, findErr := s.oauth.Find(provider)
	if findErr != nil {
		return empty, connectionhelpers.ErrOAuthUnavailable
	}
	token, err = strategy.Refresh(ctx, token)
	if err != nil {
		return empty, err
	}
	encrypted, err = s.cipher.Seal(token)
	if err != nil {
		return empty, err
	}
	if err = s.repository.UpdateCredential(ctx, connectionID, encrypted); err != nil {
		return empty, err
	}
	return token, nil
}


