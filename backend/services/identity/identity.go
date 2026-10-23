package identity

import (
	"log/slog"
	"net/http"

	coretelemetry "grovehub/libs/core/telemetry"
	"grovehub/services/identity/controller"
	"grovehub/services/identity/routes"
)

type Module interface {
	RegisterHTTP(mux *http.ServeMux)
}

type module struct {
	logger           *slog.Logger
	healthController http.Handler
}

func New(logger *slog.Logger) Module {
	return &module{
		logger:           logger,
		healthController: controller.NewHealthController(),
	}
}

func (service *module) RegisterHTTP(mux *http.ServeMux) {
	routes.RegisterHTTP(mux, service.healthController)
	service.logger.Info(
		"identity service registered",
		"service", coretelemetry.ServiceIdentity.Name,
		"serviceKey", coretelemetry.ServiceIdentity.Key,
	)
}
