package impl

import (
    client "grovehub/services/infrastructure/client/messaging"
    config "grovehub/services/infrastructure/internal/app/config/messaging"
    natsadapter "grovehub/services/infrastructure/internal/app/externalService/messaging/nats/adapters"
    messagingservice "grovehub/services/infrastructure/internal/app/services/messaging"
    serviceimpl "grovehub/services/infrastructure/internal/app/services/messaging/impl"
    "grovehub/services/infrastructure/internal/app/services/messaging/strategies"
)

// newMessaging assembles provider-neutral services, not an HTTP controller.
func (b *builder) newMessaging(store client.EventStore) (messagingservice.Service, error) {
    adapter, err := natsadapter.New(b.settings)
    if err != nil { return nil, err }
    registry := strategies.NewRegistry[strategies.Strategy]()
    registry.Register(config.ProviderNATS, adapter)
    service, err := serviceimpl.NewMessagingService(b.settings.Provider, registry, store)
    if err != nil {
        _ = adapter.Close()
        return nil, err
    }
    return service, nil
}
