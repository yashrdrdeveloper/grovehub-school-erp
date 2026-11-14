package impl

import (
    "context"

    foundationevents "grovehub/services/foundation/constants/events"
    persondto "grovehub/services/foundation/dto/person"
    eventpipelines "grovehub/services/foundation/events/person/pipelines"
    eventflows "grovehub/services/foundation/events/person/pipelines/flows"
    eventproducers "grovehub/services/foundation/events/person/producers"
)

type PersonProducer struct {
    pipeline *eventpipelines.ChangedPipeline
}

var _ eventproducers.ChangedProducer = (*PersonProducer)(nil)

func New(pipeline *eventpipelines.ChangedPipeline) eventproducers.ChangedProducer { return &PersonProducer{pipeline: pipeline} }

func (producer *PersonProducer) Enqueue(ctx context.Context, operation foundationevents.Operation, old, next *persondto.Person) error {
    return producer.pipeline.Enqueue(ctx, eventflows.Change{Operation: operation, Old: old, New: next})
}
