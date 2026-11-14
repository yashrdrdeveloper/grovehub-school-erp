package impl

import (
    "context"
    "time"

    coreoutbox "grovehub/libs/core/events/outbox"
    eventpipelines "grovehub/services/foundation/internal/app/events/producers/person/operations/changed"
    eventstrategyimpl "grovehub/services/foundation/internal/app/events/producers/person/operations/changed/strategies/impl"
    personproducer "grovehub/services/foundation/internal/app/events/producers/person"
    eventproducerimpl "grovehub/services/foundation/internal/app/events/producers/person/impl"
)

type eventRuntime struct {
    producer personproducer.ChangedProducer
    dispatcher *coreoutbox.Dispatcher
}

func (bootstrap *builder) newEvents() eventRuntime {
    store := coreoutbox.NewGormStore(bootstrap.db, "foundation.event_outbox")
    pipeline := eventpipelines.NewChangedPipeline(eventstrategyimpl.New(), store)
    return eventRuntime{
        producer: eventproducerimpl.New(pipeline),
        dispatcher: coreoutbox.NewDispatcher(store, bootstrap.messaging, bootstrap.logger),
    }
}

func (events eventRuntime) start() (context.CancelFunc, <-chan struct{}) {
    ctx, cancel := context.WithCancel(context.Background())
    done := make(chan struct{})
    go func() {
        defer close(done)
        events.dispatcher.Run(ctx, time.Second)
    }()
    return cancel, done
}
