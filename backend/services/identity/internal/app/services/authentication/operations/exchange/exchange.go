package exchange

import (
	"context"
	auth "grovehub/services/identity/internal/app/services/authentication"
)

// Operation is the use-case contract; concrete dependencies live in impl/.
type Operation interface {
	Execute(ctx context.Context, secret string, purpose string) (auth.Grant, error)
}
