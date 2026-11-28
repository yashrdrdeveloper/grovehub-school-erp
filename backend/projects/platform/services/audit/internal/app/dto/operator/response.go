package operator

import (
	"time"

	"grovehub/services/audit/internal/app/dto/producer"
)

type SearchPage struct {
	Records    []producer.OperatorActivity `json:"records"`
	PageNumber int                         `json:"pageNumber"`
	RowSize    int                         `json:"rowSize"`
	HasMore    bool                        `json:"hasMore"`
}

type ReceiptView struct {
	ExecutionID     string    `json:"executionId"`
	SchoolID        string    `json:"schoolId"`
	Operation       string    `json:"operation"`
	PreviousVersion string    `json:"previousVersion"`
	CorrelationID   *string   `json:"correlationId,omitempty"`
	OccurredAt      time.Time `json:"occurredAt"`
}

type ReceiptPage struct {
	Records    []ReceiptView `json:"records"`
	PageNumber int           `json:"pageNumber"`
	RowSize    int           `json:"rowSize"`
	HasMore    bool          `json:"hasMore"`
}
