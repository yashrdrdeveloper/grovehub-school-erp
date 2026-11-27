// Package principal carries identities established by a trusted gateway.
package principal

import "context"

type VerifiedPrincipal struct {
	TenantID string
	ActorGID string
}

type contextKey struct{}

// WithVerifiedPrincipal must only be called after authentication at a trusted
// entrypoint. Untrusted HTTP headers cannot establish a principal.
func WithVerifiedPrincipal(ctx context.Context, value VerifiedPrincipal) context.Context {
	return context.WithValue(ctx, contextKey{}, value)
}

func From(ctx context.Context) (VerifiedPrincipal, bool) {
	value, ok := ctx.Value(contextKey{}).(VerifiedPrincipal)
	return value, ok && value.TenantID != "" && value.ActorGID != ""
}
