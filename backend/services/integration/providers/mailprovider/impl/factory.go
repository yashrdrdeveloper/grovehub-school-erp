package impl

import (
	"strings"

	"grovehub/services/integration/providers"
	providerimpl "grovehub/services/integration/providers/impl"
	"grovehub/services/integration/providers/mailprovider"
)

type factory struct{}

func NewFactory() mailprovider.Factory {
	return factory{}
}

type resolver struct {
	strategies map[string]mailprovider.Strategy
}

func (r *resolver) Resolve(providerID string, capability providers.Capability) (mailprovider.Strategy, error) {
	strategy, ok := r.strategies[strings.ToUpper(strings.TrimSpace(providerID))]
	if !ok {
		return nil, providers.ErrUnknownProvider
	}
	for _, supported := range strategy.Descriptor().Capabilities {
		if supported == capability {
			return strategy, nil
		}
	}
	return nil, providers.ErrUnsupportedAction
}

func (factory) Build(strategies ...mailprovider.Strategy) (mailprovider.Runtime, error) {
	catalog := make([]providers.Provider, 0, len(strategies))
	selector := &resolver{strategies: make(map[string]mailprovider.Strategy, len(strategies))}
	for _, strategy := range strategies {
		if strategy == nil {
			return mailprovider.Runtime{}, providers.ErrUnknownProvider
		}
		descriptor := strategy.Descriptor()
		catalog = append(catalog, strategy)
		selector.strategies[strings.ToUpper(strings.TrimSpace(descriptor.ID))] = strategy
	}
	registry, err := providerimpl.NewRegistry(catalog...)
	if err != nil {
		return mailprovider.Runtime{}, err
	}
	actions := providerimpl.NewActionRegistry()
	oauth := providerimpl.NewOAuthRegistry()
	for _, strategy := range strategies {
		if err := strategy.Register(actions, oauth); err != nil {
			return mailprovider.Runtime{}, err
		}
	}
	return mailprovider.Runtime{
		Registry: registry,
		Actions:  actions,
		OAuth:    oauth,
		Resolver: selector,
	}, nil
}
