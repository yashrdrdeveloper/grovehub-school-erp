package impl

import "grovehub/launcher/internal/src/services/health"

// @Service
type service struct{ ready func() bool }

var _ health.Service = (*service)(nil)

func New(ready func() bool) health.Service {
	return &service{ready: ready}
}

func (*service) Health() health.Response {
	return health.Response{Status: "UP"}
}

func (s *service) Readiness() (health.Response, bool) {
	if s.ready == nil || !s.ready() {
		return health.Response{Status: "NOT_READY"}, false
	}
	return health.Response{Status: "READY"}, true
}
