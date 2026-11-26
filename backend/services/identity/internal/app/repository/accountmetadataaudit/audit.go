package accountmetadataaudit

import (
	"context"
	"time"
)

type Entry struct {
	EventID string
	TenantID string
	ActorID string
	AccountID string
	Version int64
	Payload []byte
	OccurredAt time.Time
}

type Writer interface {
	Insert(ctx context.Context, entry Entry) error
}
