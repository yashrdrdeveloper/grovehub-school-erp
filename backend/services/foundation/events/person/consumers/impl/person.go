package impl

import (
    "context"

    personevents "grovehub/services/foundation/events/person"
    eventconsumers "grovehub/services/foundation/events/person/consumers"
    eventhandlers "grovehub/services/foundation/events/person/handlers"
    infrastructureclient "grovehub/services/infrastructure/client"
)

type ChangedConsumer struct {
    messaging infrastructureclient.MessagingClient
    handle eventhandlers.ChangedHandler
}

var _ eventconsumers.ChangedConsumer = (*ChangedConsumer)(nil)

func New(messaging infrastructureclient.MessagingClient, handler eventhandlers.ChangedHandler) eventconsumers.ChangedConsumer {
    return &ChangedConsumer{messaging: messaging, handle: handler}
}

func (consumer *ChangedConsumer) Register(ctx context.Context) error {
    return consumer.messaging.Consume(ctx, personevents.ChangedTopic,
        func(ctx context.Context, message infrastructureclient.Message) error {
            return eventhandlers.DecodeChanged(ctx, message.Payload, consumer.handle)
        })
}
