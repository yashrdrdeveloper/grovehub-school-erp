package persons

import (
    foundationevents "grovehub/services/foundation/constants/events"
    persondto "grovehub/services/foundation/dto/person"
)

// Changed carries full public Person response snapshots. Nil Old means creation;
// nil New is reserved for deletion, should the service later support it.
// Consumers must protect contact PII and should not log raw event payloads.
// @Event topic=foundation_person_changed type=foundation.person.changed.v1 source=foundation
type Changed struct {
    Operation foundationevents.Operation `json:"operation"`
    Old *persondto.Person `json:"old"`
    New *persondto.Person `json:"new"`
}
