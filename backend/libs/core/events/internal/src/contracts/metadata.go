package contracts

import "time"

// Metadata contains domain-agnostic event identity and schema metadata.
type Metadata struct {
	EventID       string
	EventType     string
	Version       int
	AggregateID   string
	OccurredAt    time.Time
	CorrelationID string
	CausationID   string
}
