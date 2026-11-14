// Package pipelines executes Person change event rules inside the business transaction.
package changed

import (
    "context"

    coreoutbox "grovehub/libs/core/events/outbox"
    eventpipeline "grovehub/libs/core/events/pipeline"
    eventflows "grovehub/services/foundation/internal/app/events/producers/person/operations/changed/flows"
    changedstrategies "grovehub/services/foundation/internal/app/events/producers/person/operations/changed/strategies"
)

type ChangedPipeline struct {
    strategy changedstrategies.ChangedStrategy
    store coreoutbox.Store
}

func NewChangedPipeline(strategy changedstrategies.ChangedStrategy, store coreoutbox.Store) *ChangedPipeline {
    return &ChangedPipeline{strategy: strategy, store: store}
}

func (pipeline *ChangedPipeline) Rules() eventpipeline.Rules[eventflows.Change] {
    return eventpipeline.Rules[eventflows.Change]{
        ShouldEmit: pipeline.strategy.ShouldEmit,
        Build: pipeline.strategy.Build,
        Enqueue: pipeline.store.Enqueue,
    }
}

func (pipeline *ChangedPipeline) Enqueue(ctx context.Context, change eventflows.Change) error {
    return eventpipeline.Run(ctx, change, pipeline.Rules())
}
