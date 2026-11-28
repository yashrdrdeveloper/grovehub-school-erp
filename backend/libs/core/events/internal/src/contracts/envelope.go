package contracts

import "time"

// Envelope defines the transport-independent event wire contract.
// Version is the event schema version, not the aggregate's optimistic-lock version.
// CorrelationID and CausationID are omitted when unknown.
type Envelope[T any] struct {
	EventID       string    `json:"eventId"`
	EventType     string    `json:"eventType"`
	Version       int       `json:"version"`
	AggregateID   string    `json:"aggregateId"`
	OccurredAt    time.Time `json:"occurredAt"`
	CorrelationID string    `json:"correlationId,omitempty"`
	CausationID   string    `json:"causationId,omitempty"`
	Data          T         `json:"data"`
}
