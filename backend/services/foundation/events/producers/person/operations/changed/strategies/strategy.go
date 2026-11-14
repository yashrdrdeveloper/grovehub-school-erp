package strategies

import (
    "context"

    coreoutbox "grovehub/libs/core/events/outbox"
    eventflows "grovehub/services/foundation/events/producers/person/operations/changed/flows"
)

// @Strategy
type ChangedStrategy interface {
    ShouldEmit(eventflows.Change) bool
    Build(context.Context, eventflows.Change) (coreoutbox.Message, error)
}
