package bootstrap

import (
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/module"
	schoolgateway "grovehub/libs/schoolgateway"
	accountclient "grovehub/services/identity/client/account"
)

type Runtime interface {
	module.Module
	schoolgateway.AccessVerifier
	schoolgateway.ProtectedRouteAuthorizer
	AccountClient() accountclient.Client
	RouteLimits() map[string]int
}

type Bootstrap interface {
	Build() (Runtime, error)
}
