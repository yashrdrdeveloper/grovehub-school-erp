package strategies

import (
	"context"

	"grovehub/services/identity/internal/app/dto/authentication"
	auth "grovehub/services/identity/internal/app/services/authentication"
)

// LoginStrategy authenticates one credential family. Registration,
// token recovery, and account lifecycle do not belong in this interface.
// @Strategy
type LoginStrategy interface {
	Login(ctx context.Context, input authentication.Login, userAgent string) (auth.SessionResult, error)
}
