package rolemetadata

import (
	"context"

	"grovehub/services/identity/internal/app/events/producers/resource/operations/changed"
)

type Recorder interface {
	Record(ctx context.Context, before *changed.Snapshot, after changed.Snapshot) error
}
