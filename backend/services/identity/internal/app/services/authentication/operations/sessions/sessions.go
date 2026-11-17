package sessions

import (
	"context"
	authdto "grovehub/services/identity/internal/app/dto/authentication"
)

// Operation is the use-case contract; concrete dependencies live in impl/.
type Operation interface {
	Execute(ctx context.Context) (authdto.SessionsResult, error)
}
