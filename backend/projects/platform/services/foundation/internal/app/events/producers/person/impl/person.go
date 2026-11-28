package impl

import (
	"context"
	schoolgateway "grovehub/libs/platform/session"
	auditclient "grovehub/services/audit/client/producer"
	eventhelpers "grovehub/services/foundation/internal/app/events/producers/person/operations/changed/helpers"

	foundationevents "grovehub/services/foundation/internal/app/constants/events"
	persondto "grovehub/services/foundation/internal/app/dto/person"
	eventproducers "grovehub/services/foundation/internal/app/events/producers/person"
	eventpipelines "grovehub/services/foundation/internal/app/events/producers/person/operations/changed"
	eventflows "grovehub/services/foundation/internal/app/events/producers/person/operations/changed/flows"
)

type PersonProducer struct {
	pipeline *eventpipelines.ChangedPipeline
	audit    auditclient.Client
}

var _ eventproducers.ChangedProducer = (*PersonProducer)(nil)

func New(pipeline *eventpipelines.ChangedPipeline, clients ...auditclient.Client) eventproducers.ChangedProducer {
	producer := &PersonProducer{pipeline: pipeline}
	if len(clients) > 0 {
		producer.audit = clients[0]
	}
	return producer
}

func (producer *PersonProducer) Enqueue(ctx context.Context, operation foundationevents.Operation, old, next *persondto.Person) error {
	change := eventflows.Change{Operation: operation, Old: old, New: next}
	if err := producer.pipeline.Enqueue(ctx, change); err != nil {
		return err
	}
	if producer.audit == nil || !eventhelpers.HasChanged(change) {
		return nil
	}
	// Foundation persons are global. Only a trusted tenant principal can
	// establish the scope of the activity, never the person's contacts.
	actor, ok := schoolgateway.PrincipalFrom(ctx)
	if !ok || actor.TenantID == "" {
		return nil
	}
	return producer.audit.Record(ctx, auditclient.Activity{TenantID: actor.TenantID, Source: "foundation", AssetType: "PERSON", AssetID: next.GID, PersonID: next.GID, ActorID: actor.AccountID, Operation: string(operation), AssetVersion: next.Version, OccurredAt: next.UpdatedAt})
}
