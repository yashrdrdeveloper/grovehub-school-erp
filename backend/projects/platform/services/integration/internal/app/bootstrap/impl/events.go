package impl

import (
    producer "grovehub/services/integration/internal/app/events/producers/connection"
    producerimpl "grovehub/services/integration/internal/app/events/producers/connection/impl"
    changed "grovehub/services/integration/internal/app/events/producers/connection/operations/mailbox"
    changedimpl "grovehub/services/integration/internal/app/events/producers/connection/operations/mailbox/strategies/impl"
)

// newEvents follows Foundation's event-producer composition pattern.
func (b *builder) newEvents() producer.ChangedProducer {
    pipeline:=changed.NewChangedPipeline(changedimpl.New(),b.messaging)
    return producerimpl.New(pipeline)
}
