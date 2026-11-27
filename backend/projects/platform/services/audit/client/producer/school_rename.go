package producer

import "time"

// SchoolRenameReceipt records a committed mutation, not an asserted user activity.
type SchoolRenameReceipt struct{
 ExecutionID string `json:"executionId"`
 SchoolID string `json:"schoolId"`
 PreviousVersion int64 `json:"previousVersion"`
 CorrelationID *string `json:"correlationId,omitempty"`
 OccurredAt time.Time `json:"occurredAt"`
}
