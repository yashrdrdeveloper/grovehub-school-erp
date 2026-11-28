package security

import (
	"context"

	schoolgateway "grovehub/internal/session"
)

// AccessVerifier is the Identity-owned contract for validating a browser
// session. Launcher only depends on this interface, never SQL repositories.
type AccessVerifier interface {
	VerifyAccess(ctx context.Context, secret string) (schoolgateway.Principal, error)
}
