package strategies

import (
	"context"

	eventflows "grovehub/services/foundation/internal/app/events/producers/person/operations/changed/flows"
	messagingclient "grovehub/services/infrastructure/client/messaging"
)

// @Strategy
type ChangedStrategy interface {
	ShouldEmit(eventflows.Change) bool
	Build(context.Context, eventflows.Change) (messagingclient.Event, error)
}
