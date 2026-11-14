package strategies

import (
    "context"

    messagingclient "grovehub/services/infrastructure/client/messaging"
    eventflows "grovehub/services/foundation/internal/app/events/producers/person/operations/changed/flows"
)

// @Strategy
type ChangedStrategy interface {
    ShouldEmit(eventflows.Change) bool
    Build(context.Context, eventflows.Change) (messagingclient.Event, error)
}
