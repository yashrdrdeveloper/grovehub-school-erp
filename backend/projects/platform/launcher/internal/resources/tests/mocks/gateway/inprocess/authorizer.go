package inprocess

import (
	"context"

	"grovehub/libs/platform/session"
)

// @Mock
type Authorizer interface {
	Authorize(context.Context, session.Principal, string) error
}
