package entity

import "time"

// @Entity audit.operator_activities
type OperatorActivity struct {
	EventID       string `gorm:"primaryKey;type:uuid"`
	Source        string
	ActorID       string `gorm:"type:uuid"`
	ResourceID    string `gorm:"type:uuid"`
	Operation     string
	OccurredAt    time.Time
	CorrelationID *string `gorm:"type:uuid"`
}
