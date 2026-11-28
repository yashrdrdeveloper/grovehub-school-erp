package external

import (
	"context"
	"net/http"

	"grovehub/libs/platform/session"
)

// @Mock
type AccessVerifier interface {
	VerifyAccess(context.Context, string) (session.Principal, error)
}

// @Mock
type WorkloadVerifier interface {
	VerifyWorkload(context.Context, *http.Request) (context.Context, error)
}
