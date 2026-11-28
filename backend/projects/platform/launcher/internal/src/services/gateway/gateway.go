// Package gateway composes the external and internal gateway operations.
package gateway

import (
	"net/http"

	external "grovehub/launcher/internal/src/services/gateway/operations/external"
	inprocess "grovehub/launcher/internal/src/services/gateway/operations/inprocess"
	coregateway "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
	schoolgateway "grovehub/launcher/services/gateway/session"
)

type Gateway interface {
 External() http.Handler
 Internal() inprocess.Operation
}

type service struct {
 external http.Handler
 internal inprocess.Operation
}

var _ Gateway = (*service)(nil)

func (s *service) External() http.Handler {return s.external}
func (s *service) Internal() inprocess.Operation {return s.internal}

func New(mux *http.ServeMux, access schoolgateway.AccessVerifier, workload coregateway.WorkloadVerifier, origin string, authorizer inprocess.Authorizer, limits ...map[string]int) Gateway {
	return &service{
		external: external.New(mux, access, workload, origin, limits...),
		internal: inprocess.New(authorizer),
	}
}
