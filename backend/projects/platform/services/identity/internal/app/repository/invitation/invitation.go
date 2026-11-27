package invitation

import (
	"context"
	"time"

	"grovehub/services/identity/internal/app/entity"
)

type Repository interface {
	Queue(ctx context.Context, invitation entity.Invitation) error
	Claim(ctx context.Context, now time.Time) (entity.Invitation, bool, error)
	LockClaim(ctx context.Context, invitation entity.Invitation, now time.Time) (bool, error)
	Suppress(ctx context.Context, invitation entity.Invitation, reason string) error
	Delivered(ctx context.Context, invitationID string) error
	Retry(ctx context.Context, invitationID string, nextAttempt time.Time) error
}
