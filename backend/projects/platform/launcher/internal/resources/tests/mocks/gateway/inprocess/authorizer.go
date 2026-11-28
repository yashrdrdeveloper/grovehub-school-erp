package inprocess

import (
 "context"
 "grovehub/internal/session"
)

// @Mock
type Authorizer interface {
 Authorize(context.Context,session.Principal,string) error
}
