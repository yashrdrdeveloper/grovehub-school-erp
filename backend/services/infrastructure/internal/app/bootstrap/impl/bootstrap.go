// Package impl assembles Infrastructure's private adapters, messaging
// repository and delivery worker behind the public module entrypoint.
package impl

import (
    "context"
    "errors"
    "log/slog"
    "time"

    "gorm.io/gorm"

    messagingclient "grovehub/services/infrastructure/client/messaging"
    bootstrap "grovehub/services/infrastructure/internal/app/bootstrap"
    clientimpl "grovehub/services/infrastructure/internal/app/client/messaging/impl"
    config "grovehub/services/infrastructure/internal/app/config"
    natsadapters "grovehub/services/infrastructure/internal/app/externalService/nats/adapters"
    repository "grovehub/services/infrastructure/internal/app/repository/messaging/impl"
    mappers "grovehub/services/infrastructure/internal/app/mappers/messaging/impl"
    service "grovehub/services/infrastructure/internal/app/services/messaging"
    serviceimpl "grovehub/services/infrastructure/internal/app/services/messaging/impl"
    "grovehub/services/infrastructure/internal/app/services/messaging/strategies"
    validators "grovehub/services/infrastructure/internal/app/validators/messaging/impl"
)

type runtime struct {
    messagingClient messagingclient.MessagingClient
    eventStore messagingclient.EventStore
    provider interface{Close() error}
    stop context.CancelFunc
    done <-chan struct{}
}

var _ bootstrap.Runtime = (*runtime)(nil)

func New(settings config.Messaging, db *gorm.DB, logger *slog.Logger) (bootstrap.Runtime,error) {
    if db == nil { return nil, errors.New("infrastructure messaging requires the shared PostgreSQL database") }
    if logger == nil { logger = slog.Default() }
    adapter,err:=natsadapters.New(settings)
    if err!=nil{return nil,err}
    registry:=strategies.NewRegistry[strategies.Strategy]()
    registry.Register(config.ProviderNATS,strategies.NewNATSStrategy(adapter))
    provider:=serviceimpl.NewMessagingService(settings.Provider,registry)
    client:=clientimpl.NewMessagingClient(provider)
    store:=repository.New(db,validators.New(),mappers.New())
    dispatcher:=service.NewDispatcher(store,client,logger)
    ctx,cancel:=context.WithCancel(context.Background())
    done:=make(chan struct{})
    go func(){
        defer close(done)
        dispatcher.Run(ctx,time.Second)
    }()
    return &runtime{messagingClient:client,eventStore:store,provider:provider,stop:cancel,done:done},nil
}

func (module *runtime) MessagingClient() messagingclient.MessagingClient {return module.messagingClient}
func (module *runtime) EventStore() messagingclient.EventStore {return module.eventStore}

// Close drains Infrastructure's dispatcher before closing the provider.
func (module *runtime) Close() error {
    module.stop()
    <-module.done
    return module.provider.Close()
}
