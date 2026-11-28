package messaging

import (
	"context"
	"errors"
)

var ErrReplayUnavailable=errors.New("published event unavailable")

// ReplaySource is an operator-only, read-only view of retained published events.
// Callers must validate the original envelope and tenant before using it.
type ReplaySource interface {
	Published(ctx context.Context,eventID string) (Event,error)
}
