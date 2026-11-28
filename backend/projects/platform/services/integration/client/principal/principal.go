// Package principal carries identities established by a trusted gateway.
package principal

import (
	"context"

	"grovehub/services/integration/internal/app/client/principal/helpers"
)

type VerifiedPrincipal = helpers.VerifiedPrincipal

// WithVerifiedPrincipal must only be called after authentication at a trusted
// entrypoint. Untrusted HTTP headers cannot establish a principal.
func WithVerifiedPrincipal(ctx context.Context, value VerifiedPrincipal) context.Context {
	return helpers.WithVerifiedPrincipal(ctx, value)
}

func From(ctx context.Context) (VerifiedPrincipal, bool) {
	return helpers.From(ctx)
}
