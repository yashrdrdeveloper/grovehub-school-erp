package workload

import (
	"context"

	"grovehub/libs/platform/session"
)

// @Mock
type DelegatedIdentity interface {
	VerifyDelegatedSession(context.Context, string) (session.Principal, error)
	VerifyAccess(context.Context, string) (session.Principal, error)
	AuthorizeProtectedRoute(context.Context, string) error
}
