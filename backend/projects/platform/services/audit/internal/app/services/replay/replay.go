package replay

import "context"

type Request struct {
	TenantID string
	EventID string
	OperatorLabel string
	Apply bool
}

type Result struct {
	EventID string `json:"eventId"`
	AuditID string `json:"auditId"`
	Operation string `json:"operation"`
	Applied bool `json:"applied"`
	ReplayID string `json:"replayId,omitempty"`
}

type Service interface {
	Replay(ctx context.Context,input Request) (Result,error)
}
