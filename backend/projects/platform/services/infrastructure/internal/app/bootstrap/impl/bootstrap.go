// Package impl composes Infrastructure's messaging runtime like Foundation's builder.
package impl

import (
    "context"
    "errors"
    "log/slog"

    "gorm.io/gorm"

    client "grovehub/services/infrastructure/client/messaging"
    bootstrap "grovehub/services/infrastructure/internal/app/bootstrap"
    config "grovehub/services/infrastructure/internal/app/config/messaging"
    messagingservice "grovehub/services/infrastructure/internal/app/services/messaging"
)

type builder struct {
    settings config.Messaging
    db *gorm.DB
    logger *slog.Logger
}

var _ bootstrap.Bootstrap = (*builder)(nil)

func New(settings config.Messaging, db *gorm.DB, logger *slog.Logger) bootstrap.Bootstrap {
    return &builder{settings: settings, db: db, logger: logger}
}

type runtime struct {
    messagingClient client.MessagingClient
    eventStore client.EventStore
    provider messagingservice.Service
    stop context.CancelFunc
    done <-chan struct{}
}

var _ bootstrap.Runtime = (*runtime)(nil)

func (b *builder) Build() (bootstrap.Runtime, error) {
    if b.db == nil { return nil, errors.New("infrastructure messaging requires the shared PostgreSQL database") }
    if err := b.settings.Validate(); err != nil { return nil, err }
    if b.logger == nil { b.logger = slog.Default() }
    store := b.newRepository()
    service, err := b.newMessaging(store)
    if err != nil { return nil, err }
    client := b.newClient(service)
    stop, done := b.newEvents(store, client)
    return &runtime{messagingClient: client, eventStore: store, provider: service, stop: stop, done: done}, nil
}

func (module *runtime) MessagingClient() client.MessagingClient { return module.messagingClient }
func (module *runtime) EventStore() client.EventStore { return module.eventStore }

func (module *runtime) Close() error {
    module.stop()
    <-module.done
    return module.provider.Close()
}
