package helpers

import (
    foundationevents "grovehub/services/foundation/constants/events"
    eventflows "grovehub/services/foundation/events/person/pipelines/flows"
)

// HasChanged makes eligibility explicit and suppresses replay, no-op and read events.
func HasChanged(change eventflows.Change) bool {
    predicates := map[foundationevents.Operation]func(eventflows.Change) bool{
        foundationevents.Created: func(c eventflows.Change) bool { return c.Old == nil && c.New != nil },
        foundationevents.Updated: func(c eventflows.Change) bool { return c.Old != nil && c.New != nil && c.Old.Version != c.New.Version },
        foundationevents.Deactivated: func(c eventflows.Change) bool { return c.Old != nil && c.New != nil && c.Old.Version != c.New.Version },
    }
    predicate, ok := predicates[change.Operation]
    return ok && predicate(change)
}
