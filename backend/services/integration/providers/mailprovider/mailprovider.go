package mailprovider

import "grovehub/services/integration/providers"

// Strategy exposes one mail provider's catalog and optional executable
// capabilities. The factory composes strategies, not vendor-specific switches.
type Strategy interface{
 providers.Provider
 Register(actions providers.ActionRegistry,oauth providers.OAuthRegistry)error
}

// Runtime provides type-agnostic registries for consumers.
type Runtime struct{
 Registry providers.Registry
 Actions providers.ActionRegistry
 OAuth providers.OAuthRegistry
}

type Factory interface{
 Build(strategies ...Strategy)(Runtime,error)
}
