package entity

import "time"

// @Entity audit.operator_activities
type OperatorActivity struct {
	EventID string `gorm:"primaryKey;type:uuid"`
	Source string
	ActorID string `gorm:"type:uuid"`
	ResourceID string `gorm:"type:uuid"`
	Operation string
	OccurredAt time.Time
	CorrelationID *string `gorm:"type:uuid"`
}

// @Entity audit.school_mutation_receipts
type SchoolMutationReceipt struct {
	ExecutionID string `gorm:"primaryKey;type:uuid"`
	SchoolID string `gorm:"type:uuid"`
	PreviousVersion int64
	CorrelationID *string `gorm:"type:uuid"`
	OccurredAt time.Time
	Operation string
}
