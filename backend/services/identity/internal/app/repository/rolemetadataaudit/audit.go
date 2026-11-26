package rolemetadataaudit

import (
	"context"
	"time"
)

type Entry struct {
	EventID string
	TenantID string
	ActorID string
	RoleID string
	Version int64
	Operation string
	Payload []byte
	OccurredAt time.Time
}

type Writer interface {
	Insert(ctx context.Context, entry Entry) error
}
