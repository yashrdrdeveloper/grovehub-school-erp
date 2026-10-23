package foundation

import (
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	coretelemetry "grovehub/libs/core/telemetry"
	"grovehub/services/foundation/controller"
	foundationpostgres "grovehub/services/foundation/repository/postgres"
	"grovehub/services/foundation/routes"
	"grovehub/services/foundation/services"
)

type Module interface {
	RegisterHTTP(mux *http.ServeMux)
}

type module struct {
	logger           *slog.Logger
	personController *controller.PersonController
}

func New(logger *slog.Logger, pool *pgxpool.Pool) Module {
	personRepository := foundationpostgres.NewPersonRepository(pool)
	personService := services.NewPersonService(personRepository)

	return &module{
		logger:           logger,
		personController: controller.NewPersonController(personService),
	}
}

func (service *module) RegisterHTTP(mux *http.ServeMux) {
	routes.RegisterHTTP(mux, service.personController)
	service.logger.Info(
		"foundation service registered",
		"service", coretelemetry.ServiceFoundation.Name,
		"serviceKey", coretelemetry.ServiceFoundation.Key,
	)
}
