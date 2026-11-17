package revokesession

import (
	"context"
)

// Operation is the use-case contract; concrete dependencies live in impl/.
type Operation interface {
	Execute(ctx context.Context, sessionID string) error
}
