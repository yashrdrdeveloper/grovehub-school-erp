// Package gateway composes the external and internal gateway operations.
package gateway

import (
	"net/http"

	coregateway "grovehub/libs/core/gateway"
	external "grovehub/launcher/services/gateway/operations/external"
	inner "grovehub/launcher/services/gateway/operations/internal"
)

type Gateway struct {
	External http.Handler
	Internal *inner.Operation
}

func New(mux *http.ServeMux, access coregateway.AccessVerifier, workload coregateway.WorkloadVerifier, origin string, authorizer inner.Authorizer) Gateway {
	return Gateway{
		External: external.New(mux, access, workload, origin),
		Internal: inner.New(authorizer),
	}
}
