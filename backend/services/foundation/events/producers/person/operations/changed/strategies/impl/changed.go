package impl

import (
    "context"

    coreoutbox "grovehub/libs/core/events/outbox"
    eventhelpers "grovehub/services/foundation/events/producers/person/operations/changed/helpers"
    eventflows "grovehub/services/foundation/events/producers/person/operations/changed/flows"
    changedstrategies "grovehub/services/foundation/events/producers/person/operations/changed/strategies"
)

// @Strategy
type ChangedStrategy struct{}

var _ changedstrategies.ChangedStrategy = (*ChangedStrategy)(nil)

func New() changedstrategies.ChangedStrategy { return &ChangedStrategy{} }
func (strategy *ChangedStrategy) ShouldEmit(change eventflows.Change) bool { return eventhelpers.HasChanged(change) }
func (strategy *ChangedStrategy) Build(ctx context.Context, change eventflows.Change) (coreoutbox.Message, error) {
    return eventhelpers.BuildMessage(ctx, change)
}
