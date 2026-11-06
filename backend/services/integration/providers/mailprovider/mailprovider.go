package mailprovider

import "grovehub/services/integration/providers"

// Strategy describes one mail provider plus optional executable actions/OAuth.
type Strategy interface {
 providers.Provider
 Register(actions providers.ActionRegistry, oauth providers.OAuthRegistry) error
}

// Resolver selects providers by identifier and declared capability.
// Declared capability support does not imply that OAuth is configured.
type Resolver interface {
 Resolve(providerID string, capability providers.Capability) (Strategy, error)
}

type Runtime struct {
 Registry providers.Registry
 Actions providers.ActionRegistry
 OAuth providers.OAuthRegistry
 Resolver Resolver
}

type Factory interface {
 Build(strategies ...Strategy) (Runtime, error)
}
