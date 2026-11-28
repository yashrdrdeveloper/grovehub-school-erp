// Package pipelines executes Person change event rules inside the business transaction.
package changed

import (
	"context"

	coreevents "grovehub/libs/core/events"
	eventflows "grovehub/services/foundation/internal/app/events/producers/person/operations/changed/flows"
	changedstrategies "grovehub/services/foundation/internal/app/events/producers/person/operations/changed/strategies"
	messagingclient "grovehub/services/infrastructure/client/messaging"
)

type ChangedPipeline struct {
	strategy changedstrategies.ChangedStrategy
	store    messagingclient.EventStore
}

func NewChangedPipeline(strategy changedstrategies.ChangedStrategy, store messagingclient.EventStore) *ChangedPipeline {
	return &ChangedPipeline{strategy: strategy, store: store}
}

func (pipeline *ChangedPipeline) Rules() coreevents.Rules[eventflows.Change, messagingclient.Event] {
	return coreevents.Rules[eventflows.Change, messagingclient.Event]{
		ShouldEmit: pipeline.strategy.ShouldEmit,
		Build:      pipeline.strategy.Build,
		Enqueue:    pipeline.store.Enqueue,
	}
}

func (pipeline *ChangedPipeline) Enqueue(ctx context.Context, change eventflows.Change) error {
	return coreevents.Run(ctx, change, pipeline.Rules())
}
