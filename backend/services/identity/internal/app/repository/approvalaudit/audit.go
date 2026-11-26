package approvalaudit

import (
	"context"
	"time"
)

type Entry struct {
	EventID string
	TenantID string
	RequestID string
	ActorID string
	Version int64
	Operation string
	Payload []byte
	OccurredAt time.Time
}

type Writer interface {
	Insert(ctx context.Context, entry Entry) error
}
