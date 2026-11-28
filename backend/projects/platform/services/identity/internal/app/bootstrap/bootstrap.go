package bootstrap

import (
	"context"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/module"
	schoolgateway "grovehub/libs/platform/session"
	accountclient "grovehub/services/identity/client/account"
)

type Runtime interface {
	module.Module
	schoolgateway.AccessVerifier
	VerifyDelegatedSession(context.Context, string) (schoolgateway.Principal, error)
	schoolgateway.ProtectedRouteAuthorizer
	AccountClient() accountclient.Client
	RouteLimits() map[string]int
}

type Bootstrap interface {
	Build() (Runtime, error)
}
