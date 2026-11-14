package impl

import (
    "context"

    foundationevents "grovehub/services/foundation/internal/app/constants/events"
    persondto "grovehub/services/foundation/internal/app/dto/person"
    eventpipelines "grovehub/services/foundation/internal/app/events/producers/person/operations/changed"
    eventflows "grovehub/services/foundation/internal/app/events/producers/person/operations/changed/flows"
    eventproducers "grovehub/services/foundation/internal/app/events/producers/person"
)

type PersonProducer struct {
    pipeline *eventpipelines.ChangedPipeline
}

var _ eventproducers.ChangedProducer = (*PersonProducer)(nil)

func New(pipeline *eventpipelines.ChangedPipeline) eventproducers.ChangedProducer { return &PersonProducer{pipeline: pipeline} }

func (producer *PersonProducer) Enqueue(ctx context.Context, operation foundationevents.Operation, old, next *persondto.Person) error {
    return producer.pipeline.Enqueue(ctx, eventflows.Change{Operation: operation, Old: old, New: next})
}
