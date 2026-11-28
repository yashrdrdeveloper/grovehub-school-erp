package impl

import (
	"net/http"

	coregateway "grovehub/libs/core/gateway"
	"grovehub/launcher/internal/src/services/gateway"
	external "grovehub/launcher/internal/src/services/gateway/operations/external/impl"
	"grovehub/launcher/internal/src/services/gateway/operations/inprocess"
	inprocessimpl "grovehub/launcher/internal/src/services/gateway/operations/inprocess/impl"
	schoolgateway "grovehub/libs/platform/session"
)

// @Service
type service struct {
	external http.Handler
	internal inprocess.Operation
}

var _ gateway.Gateway = (*service)(nil)

func (s *service) External() http.Handler {
	return s.external
}

func (s *service) Internal() inprocess.Operation {
	return s.internal
}

func New(mux *http.ServeMux, access schoolgateway.AccessVerifier, workload coregateway.WorkloadVerifier, origin string, authorizer inprocess.Authorizer, limits ...map[string]int) gateway.Gateway {
	return &service{
		external: external.New(mux, access, workload, origin, limits...),
		internal: inprocessimpl.New(authorizer),
	}
}
