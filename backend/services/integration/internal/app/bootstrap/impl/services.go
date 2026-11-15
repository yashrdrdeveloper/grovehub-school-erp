package impl

import (
    validatorimpl "grovehub/services/integration/internal/app/validators/connection/impl"
    connectionservice "grovehub/services/integration/internal/app/services/connection"
    connectionimpl "grovehub/services/integration/internal/app/services/connection/impl"
    providerservice "grovehub/services/integration/internal/app/services/provider"
    providerimpl "grovehub/services/integration/internal/app/services/provider/impl"
)

func (b *builder) newServices(providers providerDependencies) (connectionservice.ConnectionService, providerservice.ProviderService, error) {
    connections,watches:=b.newRepository()
    service,err:=connectionimpl.NewConnectionService(
        connections,watches,
        providers.registry,providers.actions,providers.oauth,
        providers.gmail,providers.outlook,providers.key,
        providers.push.Topic,b.messaging,validatorimpl.NewConnectionValidator(),
    )
    if err!=nil{return nil,nil,err}
    providerService:=providerimpl.NewProviderService(providers.registry,providers.actions,providers.oauth)
    return service,providerService,nil
}
