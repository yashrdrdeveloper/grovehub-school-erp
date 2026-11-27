package approval

import (
	"context"

	"grovehub/services/identity/internal/app/entity"
)

// Recorder persists an immutable audit event and enqueues its delivery using
// the caller's transaction. It must not publish to the broker directly.
type Recorder interface {
	Record(ctx context.Context, operation string, before *entity.ApprovalRequest, after entity.ApprovalRequest) error
}
