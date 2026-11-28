package principal

import (
	"context"
	"grovehub/launcher/session/contracts"
)

type principalKey struct{}

func WithPrincipal(ctx context.Context, principal contracts.Principal) context.Context {
	return context.WithValue(ctx, principalKey{}, principal)
}

func PrincipalFrom(ctx context.Context) (contracts.Principal, bool) {
	principal, ok := ctx.Value(principalKey{}).(contracts.Principal)
	return principal, ok && principal.AccountID != "" && principal.SessionID != ""
}

