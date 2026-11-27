package mailprovider

import (
    providers "grovehub/services/integration/internal/app/providers"
    contracts "grovehub/services/integration/internal/app/providers/contracts"
    dto "grovehub/services/integration/internal/app/providers/dto"
    mailcontracts "grovehub/services/integration/internal/app/providers/mailprovider/contracts"
)

// Strategy describes a mail capability adapter; provider discovery stays generic.
type Strategy interface {
    providers.Provider
    Register(actions contracts.ActionRegistry, oauth mailcontracts.OAuthRegistry) error
}

type Resolver interface {
    Resolve(providerID string, capability dto.Capability) (Strategy,error)
}

type Runtime struct {
    Registry contracts.Registry
    Actions contracts.ActionRegistry
    OAuth mailcontracts.OAuthRegistry
    Resolver Resolver
}

type Factory interface {
    Build(strategies ...Strategy) (Runtime,error)
}
