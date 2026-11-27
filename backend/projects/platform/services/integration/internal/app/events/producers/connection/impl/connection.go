package impl

import (
    "context"

    producer "grovehub/services/integration/internal/app/events/producers/connection"
    mailbox "grovehub/services/integration/internal/app/events/producers/connection/operations/mailbox"
)

type ConnectionProducer struct { pipeline *mailbox.ChangedPipeline }

var _ producer.ChangedProducer = (*ConnectionProducer)(nil)

func New(pipeline *mailbox.ChangedPipeline) producer.ChangedProducer {
    return &ConnectionProducer{pipeline:pipeline}
}

func (producer *ConnectionProducer) PublishChanged(ctx context.Context, payload mailbox.Payload) error {
    return producer.pipeline.Publish(ctx,payload)
}
