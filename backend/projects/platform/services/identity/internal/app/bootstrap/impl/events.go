package impl

import (
	producer "grovehub/services/identity/internal/app/events/producers/account"
	producerimpl "grovehub/services/identity/internal/app/events/producers/account/impl"
)

func (b *builder) newEvents() producer.ChangedProducer {
	return producerimpl.New(b.events)
}
