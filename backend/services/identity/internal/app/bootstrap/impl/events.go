package impl
import(
 producer "grovehub/services/identity/internal/app/events/producers/account"
 impl "grovehub/services/identity/internal/app/events/producers/account/impl"
)
func(b *builder)newEvents()producer.ChangedProducer{return impl.New(b.events)}
