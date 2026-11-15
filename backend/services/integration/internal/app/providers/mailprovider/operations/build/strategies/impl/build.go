package impl

import (
	"strings"

	"grovehub/services/integration/internal/app/providers"
	register "grovehub/services/integration/internal/app/providers/operations/register"
    execute "grovehub/services/integration/internal/app/providers/operations/execute"
    oauth "grovehub/services/integration/internal/app/providers/mailprovider/operations/oauth"
    dto "grovehub/services/integration/internal/app/providers/dto"
    helpers "grovehub/services/integration/internal/app/providers/helpers"
	"grovehub/services/integration/internal/app/providers/mailprovider"
)

type factory struct{}

func NewFactory() mailprovider.Factory {
	return factory{}
}

type resolver struct {
	strategies map[string]mailprovider.Strategy
}

func (r *resolver) Resolve(providerID string, capability dto.Capability) (mailprovider.Strategy, error) {
	strategy, ok := r.strategies[strings.ToUpper(strings.TrimSpace(providerID))]
	if !ok {
		return nil, helpers.ErrUnknownProvider
	}
	for _, supported := range strategy.Descriptor().Capabilities {
		if supported == capability {
			return strategy, nil
		}
	}
	return nil, helpers.ErrUnsupportedAction
}

func (factory) Build(strategies ...mailprovider.Strategy) (mailprovider.Runtime, error) {
	catalog := make([]providers.Provider, 0, len(strategies))
	selector := &resolver{strategies: make(map[string]mailprovider.Strategy, len(strategies))}
	for _, strategy := range strategies {
		if strategy == nil {
			return mailprovider.Runtime{}, helpers.ErrUnknownProvider
		}
		descriptor := strategy.Descriptor()
		catalog = append(catalog, strategy)
		selector.strategies[strings.ToUpper(strings.TrimSpace(descriptor.ID))] = strategy
	}
	registry, err := register.NewRegistry(catalog...)
	if err != nil {
		return mailprovider.Runtime{}, err
	}
	actions := execute.NewActionRegistry()
	oauth := oauth.NewRegistry()
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
