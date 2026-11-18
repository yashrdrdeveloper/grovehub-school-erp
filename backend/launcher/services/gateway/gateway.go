// Package gateway composes the external and internal gateway operations.
package gateway

import (
	"net/http"

	coregateway "grovehub/libs/core/gateway"
	external "grovehub/launcher/services/gateway/operations/external"
	inprocess "grovehub/launcher/services/gateway/operations/inprocess"
)

type Gateway struct {
	External http.Handler
	Internal *inprocess.Operation
}

func New(mux *http.ServeMux, access coregateway.AccessVerifier, workload coregateway.WorkloadVerifier, origin string, authorizer inprocess.Authorizer, limits ...map[string]int) Gateway {
	return Gateway{
		External: external.New(mux, access, workload, origin, limits...),
		Internal: inprocess.New(authorizer),
	}
}
