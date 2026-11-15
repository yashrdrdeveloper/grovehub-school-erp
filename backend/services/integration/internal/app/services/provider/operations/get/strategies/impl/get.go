package impl

import (
    "context"
    dto "grovehub/services/integration/internal/app/dto/provider"
    "grovehub/services/integration/internal/app/providers"
    strategies "grovehub/services/integration/internal/app/services/provider/operations/get/strategies"
)

type getStrategy struct {
    registry providers.Registry
    actions providers.ActionRegistry
    oauth providers.OAuthRegistry
}
var _ strategies.Strategy = (*getStrategy)(nil)
func New(registry providers.Registry, actions providers.ActionRegistry, oauth providers.OAuthRegistry) strategies.Strategy {
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
