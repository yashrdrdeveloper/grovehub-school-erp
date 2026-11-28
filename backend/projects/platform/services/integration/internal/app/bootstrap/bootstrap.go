// Package bootstrap defines Integration's composition boundary.
package bootstrap

import (
	"net/http"

	coregateway "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
	integrationclient "grovehub/services/integration/client/connection"
)

// Runtime owns Integration's controllers, typed client, verifier, and workers.
type Runtime interface {
	RegisterHTTP(*http.ServeMux)
	Client() integrationclient.Client
	coregateway.WorkloadVerifier
	LocalOnly() bool
	Close()
}

type Bootstrap interface {
	Build() (Runtime, error)
}
