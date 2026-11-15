package impl

import (
    services "grovehub/services/integration/internal/app/services"
    serviceimpl "grovehub/services/integration/internal/app/services/impl"
)

func (b *builder) newServices(providers providerDependencies) (services.ConnectionService,services.ProviderService,error) {
    connections,watches:=b.newRepository()
    service,err:=serviceimpl.NewConnectionService(
        connections,watches,
        providers.registry,providers.actions,providers.oauth,
        providers.gmail,providers.outlook,providers.key,
        providers.push.Topic,b.messaging,
    )
    if err!=nil{return nil,nil,err}
    providerService:=serviceimpl.NewProviderService(providers.registry,providers.actions,providers.oauth)
    return service,providerService,nil
}
