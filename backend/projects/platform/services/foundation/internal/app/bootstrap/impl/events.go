package impl

import (
	"grovehub/services/audit"
	personproducer "grovehub/services/foundation/internal/app/events/producers/person"
	producerimpl "grovehub/services/foundation/internal/app/events/producers/person/impl"
	changed "grovehub/services/foundation/internal/app/events/producers/person/operations/changed"
	changedimpl "grovehub/services/foundation/internal/app/events/producers/person/operations/changed/strategies/impl"
)

// Foundation enqueues atomically but never runs the shared messaging worker.
func (bootstrap *builder) newEvents() personproducer.ChangedProducer {
	pipeline := changed.NewChangedPipeline(changedimpl.New(), bootstrap.events)
	return producerimpl.New(pipeline, audit.NewProducer(bootstrap.events))
}
