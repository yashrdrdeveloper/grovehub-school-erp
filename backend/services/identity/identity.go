package identity

import (
	"log/slog"
	"net/http"

	identityhealth "github.com/yashrdrdeveloper/grovehub-school-erp/backend/services/identity/launcher/health"
	identityroutes "github.com/yashrdrdeveloper/grovehub-school-erp/backend/services/identity/launcher/routes"
)

type Module interface {
	RegisterHTTP(mux *http.ServeMux)
}

type module struct {
	logger        *slog.Logger
	healthHandler http.Handler
}

func New(logger *slog.Logger) Module {
	return &module{
		logger:        logger,
		healthHandler: identityhealth.NewHandler(),
	}
}

func (service *module) RegisterHTTP(mux *http.ServeMux) {
	identityroutes.RegisterHTTP(mux, service.healthHandler)
	service.logger.Info("identity service registered", "service", "identity")
}
