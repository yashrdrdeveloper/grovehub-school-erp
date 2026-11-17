package security

import (
	"context"

	coregateway "grovehub/libs/core/gateway"
)

// AccessVerifier is the Identity-owned contract for validating a browser
// session. Launcher only depends on this interface, never SQL repositories.
type AccessVerifier interface {
	VerifyAccess(ctx context.Context, secret string) (coregateway.Principal, error)
}
