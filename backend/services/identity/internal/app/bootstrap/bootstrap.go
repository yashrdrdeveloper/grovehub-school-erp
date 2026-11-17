package bootstrap

import (
	coregateway "grovehub/libs/core/gateway"
	"grovehub/libs/core/module"
	accountclient "grovehub/services/identity/client/account"
)

type Runtime interface {
	module.Module
	coregateway.AccessVerifier
	AccountClient() accountclient.Client
}

type Bootstrap interface {
	Build() (Runtime, error)
}
