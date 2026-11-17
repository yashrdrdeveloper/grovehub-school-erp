// Package gateway contains transport-neutral verified identity contracts.
// HTTP policies and route dispatch belong to launcher/services/gateway.
package gateway

import (
	"context"
	"net/http"
	"time"
)

type Tier string

const (
	Public Tier = "public"
	Protected Tier = "protected"
	Restricted Tier = "restricted"
)

// Principal is created by an authenticated adapter, never by request headers.
type Principal struct {
	AccountID string
	PersonGID string
	TenantID string
	SessionID string
	Admin bool
	AuthenticatedAt time.Time
}

type AccessVerifier interface {
	VerifyAccess(ctx context.Context, accessSecret string) (Principal, error)
}

type WorkloadVerifier interface {
	VerifyWorkload(ctx context.Context, request *http.Request) (context.Context, error)
}

type principalKey struct{}

func WithPrincipal(ctx context.Context, principal Principal) context.Context {
	return context.WithValue(ctx, principalKey{}, principal)
}

func PrincipalFrom(ctx context.Context) (Principal, bool) {
	principal, ok := ctx.Value(principalKey{}).(Principal)
	return principal, ok && principal.AccountID != "" && principal.SessionID != ""
}

const (
	AccessCookie = "__Host-gh_access"
	RefreshCookie = "__Secure-gh_refresh"
	CSRFName = "gh_csrf"
	ResetCookie = "__Secure-gh_reset"
	SetupCookie = "__Secure-gh_setup"
)
