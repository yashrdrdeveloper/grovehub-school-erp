// Package schoolgateway owns the school ERP authenticated session contract.
package schoolgateway

import (
	"context"
	"time"
)

// Principal is created by an authenticated adapter, never by request headers.
type Principal struct {
	AccountID       string
	PersonGID       string
	TenantID        string
	SessionID       string
	SessionGeneration int64
	Admin           bool
	AuthenticatedAt time.Time
}

type AccessVerifier interface {
	VerifyAccess(ctx context.Context, accessSecret string) (Principal, error)
}

type principalKey struct{}

func WithPrincipal(ctx context.Context, principal Principal) context.Context {
	return context.WithValue(ctx, principalKey{}, principal)
}

func PrincipalFrom(ctx context.Context) (Principal, bool) {
	principal, ok := ctx.Value(principalKey{}).(Principal)
	return principal, ok && principal.AccountID != "" && principal.SessionID != ""
}

// ProtectedRouteAuthorizer evaluates a registered Identity route using the verified principal.
type ProtectedRouteAuthorizer interface {
	AuthorizeProtectedRoute(ctx context.Context, pattern string) error
}
