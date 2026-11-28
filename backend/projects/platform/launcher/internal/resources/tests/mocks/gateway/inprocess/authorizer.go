package inprocess

import (
 "context"
 "grovehub/launcher/session"
)

// @Mock
type Authorizer interface {
 Authorize(context.Context,session.Principal,string) error
}
