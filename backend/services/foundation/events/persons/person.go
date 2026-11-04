package persons

import (
    foundationevents "grovehub/services/foundation/constants/events"
    foundationresponse "grovehub/services/foundation/dto/response"
)

const (
    ChangedTopic = "foundation_person_changed"
    ChangedType = "foundation.person.changed.v1"
    Source = "foundation"
)

// Changed carries full public Person response snapshots. Nil Old means creation;
// nil New is reserved for deletion, should the service later support it.
// Consumers must protect contact PII and should not log raw event payloads.
type Changed struct {
    Operation foundationevents.Operation `json:"operation"`
    Old *foundationresponse.Person `json:"old"`
    New *foundationresponse.Person `json:"new"`
}
