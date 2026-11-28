package impl

import (
	"net/http"

	corehttp "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http"
	"grovehub/launcher/internal/src/controller"
	"grovehub/launcher/internal/src/services/health"
	healthimpl "grovehub/launcher/internal/src/services/health/impl"
)

type healthController struct{ service health.Service }

var _ controller.HealthController = (*healthController)(nil)

func New(ready func() bool) controller.HealthController {
	return &healthController{service: healthimpl.New(ready)}
}

func (c *healthController) Health(w http.ResponseWriter, _ *http.Request) {
	_ = corehttp.WriteSuccess(w, http.StatusOK, c.service.Health())
}

func (c *healthController) Readiness(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	value, ready := c.service.Readiness()
	status := http.StatusServiceUnavailable
	if ready {
		status = http.StatusOK
	}

	_ = corehttp.WriteSuccess(w, status, value)
}
