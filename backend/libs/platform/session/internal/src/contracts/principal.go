// Package contracts defines Platform gateway security boundaries without runtime dependencies.
package contracts

import (
	"context"
	"time"
)

// Principal is created by an authenticated adapter, never by request headers.
type Principal struct {
	AccountID         string
	PersonGID         string
	TenantID          string
	SessionID         string
	SessionGeneration int64
	Admin             bool
	AuthenticatedAt   time.Time
}

type AccessVerifier interface {
	VerifyAccess(ctx context.Context, accessSecret string) (Principal, error)
}

type ProtectedRouteAuthorizer interface {
	AuthorizeProtectedRoute(ctx context.Context, pattern string) error
}
