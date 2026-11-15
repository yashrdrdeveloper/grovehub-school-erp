package impl

import (
    providerdto "grovehub/services/integration/internal/app/dto/provider"
	"context"
	"grovehub/services/integration/internal/app/providers"
	services "grovehub/services/integration/internal/app/services/provider"
)

type providerService struct {
	registry providers.Registry
	actions  providers.ActionRegistry
	oauth    providers.OAuthRegistry
}

func NewProviderService(registry providers.Registry, actions providers.ActionRegistry, oauth providers.OAuthRegistry) services.ProviderService {
	return &providerService{registry: registry, actions: actions, oauth: oauth}
}

// Provider discovery is independent of connection lifecycle operations.
func (s *providerService) Providers(_ context.Context) []providers.Descriptor {
	return s.registry.List()
}

func (s *providerService) Provider(
	_ context.Context,
	id string,
) (providerdto.ProviderDetails, error) {
	provider, err := s.registry.Find(id)
	if err != nil {
		return providerdto.ProviderDetails{}, err
	}

	descriptor := provider.Descriptor()
	_, oauthErr := s.oauth.Find(descriptor.ID)

	return providerdto.ProviderDetails{
		Provider:         descriptor,
		AvailableActions: s.actions.Actions(descriptor.ID),
		OAuthConfigured:  oauthErr == nil,
	}, nil
}
