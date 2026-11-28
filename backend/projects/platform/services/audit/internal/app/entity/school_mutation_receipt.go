package entity

import "time"

// @Entity audit.school_mutation_receipts
type SchoolMutationReceipt struct {
	ExecutionID string `gorm:"primaryKey;type:uuid"`
	SchoolID string `gorm:"type:uuid"`
	PreviousVersion int64
	CorrelationID *string `gorm:"type:uuid"`
	OccurredAt time.Time
	Operation string
}
