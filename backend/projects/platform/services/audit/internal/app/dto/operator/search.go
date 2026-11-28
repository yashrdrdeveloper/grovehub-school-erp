package operator

import (
	"time"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/request"
	"grovehub/services/audit/internal/app/dto/producer"
)

type SearchInput struct {
	PageNumber int `json:"pageNumber"`
	RowSize int `json:"rowSize"`
}

type SearchPage struct {
	Records []producer.OperatorActivity `json:"records"`
	PageNumber int `json:"pageNumber"`
	RowSize int `json:"rowSize"`
	HasMore bool `json:"hasMore"`
}

type ReceiptView struct {
	ExecutionID string `json:"executionId"`
	SchoolID string `json:"schoolId"`
	Operation string `json:"operation"`
	PreviousVersion string `json:"previousVersion"`
	CorrelationID *string `json:"correlationId,omitempty"`
	OccurredAt time.Time `json:"occurredAt"`
}

type ReceiptPage struct {
	Records []ReceiptView `json:"records"`
	PageNumber int `json:"pageNumber"`
	RowSize int `json:"rowSize"`
	HasMore bool `json:"hasMore"`
}

type ReceiptSearchInput struct {
	SearchInput
	SchoolID string `json:"schoolId,omitempty"`
	Operation string `json:"operation,omitempty"`
}

type SearchRequest struct {
	request.Request
	RequestBody SearchInput
}

type ReceiptSearchRequest struct {
	request.Request
	RequestBody ReceiptSearchInput
}
