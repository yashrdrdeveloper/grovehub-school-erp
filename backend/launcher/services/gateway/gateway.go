// Package gateway composes the external and internal gateway operations.
package gateway

import (
	"net/http"

	external "grovehub/launcher/services/gateway/operations/external"
	inprocess "grovehub/launcher/services/gateway/operations/inprocess"
	coregateway "grovehub/libs/core/gateway"
	schoolgateway "grovehub/libs/schoolgateway"
)

type Gateway struct {
	External http.Handler
	Internal *inprocess.Operation
}

func New(mux *http.ServeMux, access schoolgateway.AccessVerifier, workload coregateway.WorkloadVerifier, origin string, authorizer inprocess.Authorizer, limits ...map[string]int) Gateway {
	return Gateway{
		External: external.New(mux, access, workload, origin, limits...),
		Internal: inprocess.New(authorizer),
	}
}
