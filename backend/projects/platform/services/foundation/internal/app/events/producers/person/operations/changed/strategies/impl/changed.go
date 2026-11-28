package impl

import (
	"context"

	eventflows "grovehub/services/foundation/internal/app/events/producers/person/operations/changed/flows"
	eventhelpers "grovehub/services/foundation/internal/app/events/producers/person/operations/changed/helpers"
	changedstrategies "grovehub/services/foundation/internal/app/events/producers/person/operations/changed/strategies"
	messagingclient "grovehub/services/infrastructure/client/messaging"
)

// @Strategy
type ChangedStrategy struct{}

var _ changedstrategies.ChangedStrategy = (*ChangedStrategy)(nil)

func New() changedstrategies.ChangedStrategy { return &ChangedStrategy{} }
func (strategy *ChangedStrategy) ShouldEmit(change eventflows.Change) bool {
	return eventhelpers.HasChanged(change)
}
func (strategy *ChangedStrategy) Build(ctx context.Context, change eventflows.Change) (messagingclient.Event, error) {
	return eventhelpers.BuildMessage(ctx, change)
}
