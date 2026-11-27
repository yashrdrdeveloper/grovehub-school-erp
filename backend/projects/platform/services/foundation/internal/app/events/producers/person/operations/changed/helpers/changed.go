package helpers

import (
    foundationevents "grovehub/services/foundation/internal/app/constants/events"
    eventflows "grovehub/services/foundation/internal/app/events/producers/person/operations/changed/flows"
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
