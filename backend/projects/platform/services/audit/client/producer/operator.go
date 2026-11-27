package producer

import "time"

// OperatorActivity is an occurrence asserted by the independently authenticated
// Internal Tool workload. It is not a School identity or tenant activity.
type OperatorActivity struct {
 EventID string `json:"eventId"`
 ActorID string `json:"actorId"`
 ResourceID string `json:"resourceId"`
 Operation string `json:"operation"`
 OccurredAt time.Time `json:"occurredAt"`
}
