package resend

import (
	"context"
	)

// Operation is the stable use-case contract; implementation lives in impl/.
type Operation interface {
	Execute(ctx context.Context, accountID string) error
}
