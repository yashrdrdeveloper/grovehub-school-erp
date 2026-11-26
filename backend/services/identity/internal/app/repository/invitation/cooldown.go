package invitation

import (
	"context"
	"time"
)

type CooldownReader interface {
	RecentSetup(ctx context.Context, accountID string, since time.Time) (bool, error)
}
