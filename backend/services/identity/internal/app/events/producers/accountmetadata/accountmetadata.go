package accountmetadata

import (
	"context"

	"grovehub/services/identity/internal/app/entity"
)

type Recorder interface {
	Record(ctx context.Context, before entity.Account, after entity.Account) error
}
