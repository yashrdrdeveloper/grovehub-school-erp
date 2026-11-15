package mailbox

import (
    "context"
    "errors"

    messaging "grovehub/services/infrastructure/client/messaging"
    strategies "grovehub/services/integration/internal/app/events/producers/connection/operations/mailbox/strategies"
)

// ChangedPipeline separates event serialization from immediate transport.
// This notification is not a business transaction's event-outbox operation.
type ChangedPipeline struct {
    strategy strategies.ChangedStrategy
    publisher messaging.MessagingClient
}

func NewChangedPipeline(strategy strategies.ChangedStrategy, publisher messaging.MessagingClient) *ChangedPipeline {
    return &ChangedPipeline{strategy:strategy,publisher:publisher}
}

func (pipeline *ChangedPipeline) Publish(ctx context.Context, payload Payload) error {
    if pipeline==nil || pipeline.strategy==nil || pipeline.publisher==nil {return errors.New("integration mailbox publisher is unavailable")}
    topic,data,metadata,err:=pipeline.strategy.Build(ctx,payload)
    if err!=nil {return err}
    return pipeline.publisher.Publish(ctx,topic,data,metadata)
}
