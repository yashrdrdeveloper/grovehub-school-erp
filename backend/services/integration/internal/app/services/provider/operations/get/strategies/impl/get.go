package impl

import (
    providercontracts "grovehub/services/integration/internal/app/providers/contracts"
    mailcontracts "grovehub/services/integration/internal/app/providers/mailprovider/contracts"
    "context"
    dto "grovehub/services/integration/internal/app/dto/provider"
    strategies "grovehub/services/integration/internal/app/services/provider/operations/get/strategies"
)

type getStrategy struct {
    registry providercontracts.Registry
    actions providercontracts.ActionRegistry
    oauth mailcontracts.OAuthRegistry
}

var _ strategies.Strategy = (*getStrategy)(nil)

func New(registry providercontracts.Registry, actions providercontracts.ActionRegistry, oauth mailcontracts.OAuthRegistry) strategies.Strategy {
    return &getStrategy{registry:registry,actions:actions,oauth:oauth}
}

func (s *getStrategy) Execute(_ context.Context, id string) (dto.ProviderDetails,error) {
    provider,err:=s.registry.Find(id)
    if err!=nil{return dto.ProviderDetails{},err}
    descriptor:=provider.Descriptor()
    _,oauthErr:=s.oauth.Find(descriptor.ID)
    return dto.ProviderDetails{
        Provider:descriptor,
        AvailableActions:s.actions.Actions(descriptor.ID),
        OAuthConfigured:oauthErr==nil,
    },nil
}
