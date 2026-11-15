package impl

import (
	contracts "grovehub/services/integration/internal/app/providers/mailprovider/contracts"
    helpers "grovehub/services/integration/internal/app/providers/helpers"
	"strings"
)

type oauthRegistry struct {
	strategies map[string]contracts.OAuthStrategy
}

func NewOAuthRegistry() contracts.OAuthRegistry {
	return &oauthRegistry{strategies: make(map[string]contracts.OAuthStrategy)}
}

func (r *oauthRegistry) Register(provider string, strategy contracts.OAuthStrategy) error {
	provider = strings.ToUpper(strings.TrimSpace(provider))
	if provider == "" || strategy == nil {
		return helpers.ErrInvalidActionInput
	}
	if _, exists := r.strategies[provider]; exists {
		return helpers.ErrDuplicateProvider
	}
	r.strategies[provider] = strategy
	return nil
}

func (r *oauthRegistry) Find(provider string) (contracts.OAuthStrategy, error) {
	strategy, ok := r.strategies[strings.ToUpper(strings.TrimSpace(provider))]
	if !ok {
		return nil, helpers.ErrProviderUnavailable
	}
	return strategy, nil
}
