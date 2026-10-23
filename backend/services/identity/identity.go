package identity

import (
	"log/slog"
	"net/http"

	coretelemetry "grovehub/libs/core/telemetry"
	foundationclient "grovehub/services/foundation/client"
	"grovehub/services/identity/controller"
	"grovehub/services/identity/routes"
)

type Module interface {
	RegisterHTTP(mux *http.ServeMux)
}

type module struct {
	logger           *slog.Logger
	healthController http.Handler
	foundationClient foundationclient.Client
}

func New(logger *slog.Logger, foundationClient foundationclient.Client) Module {
	return &module{
		logger:           logger,
		healthController: controller.NewHealthController(),
		foundationClient: foundationClient,
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
