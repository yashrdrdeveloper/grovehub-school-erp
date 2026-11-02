package events

import "time"

type Envelope[T any] struct {
	EventID    string    `json:"eventId"`
	EventType  string    `json:"eventType"`
	Source     string    `json:"source"`
	OccurredAt time.Time `json:"occurredAt"`
	Data       T         `json:"data"`
}
