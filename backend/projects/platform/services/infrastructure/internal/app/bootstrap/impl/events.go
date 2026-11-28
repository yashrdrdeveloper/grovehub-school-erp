package impl

import (
    "context"
    "time"

    client "grovehub/services/infrastructure/client/messaging"
    repository "grovehub/services/infrastructure/internal/app/repository/messaging"
    dispatch "grovehub/services/infrastructure/internal/app/services/messaging/operations/dispatch"
)

// Events owns the dispatcher lifecycle, not the HTTP client or publisher.
func (b *builder) newEvents(store repository.Repository, publisher client.MessagingClient) (context.CancelFunc, <-chan struct{}) {
    dispatcher := dispatch.NewDispatcher(store, publisher, b.logger)
    ctx, cancel := context.WithCancel(context.Background())
    done := make(chan struct{})
    go func() {
        defer close(done)
        dispatcher.Run(ctx, time.Second)
    }()
    return cancel, done
}
