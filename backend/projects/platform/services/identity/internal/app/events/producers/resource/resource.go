package resource

import (
	"context"
	"time"

	"grovehub/services/identity/internal/app/entity"
	"grovehub/services/identity/internal/app/events/producers/resource/operations/changed"
)

type Recorder interface {
	Record(ctx context.Context, request entity.ApprovalRequest, before changed.Snapshot, after changed.Snapshot, at time.Time) error
}
