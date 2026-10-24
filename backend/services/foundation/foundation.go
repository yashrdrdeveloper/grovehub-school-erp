package foundation

import (
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	coretelemetry "grovehub/libs/core/telemetry"
	foundationclient "grovehub/services/foundation/client"
	clientimpl "grovehub/services/foundation/client/impl"
	"grovehub/services/foundation/controller"
	controllerimpl "grovehub/services/foundation/controller/impl"
	mapperimpl "grovehub/services/foundation/mappers/impl"
	repositoryimpl "grovehub/services/foundation/repository/impl"
	"grovehub/services/foundation/routes"
	serviceimpl "grovehub/services/foundation/services/impl"
	validatorimpl "grovehub/services/foundation/validators/impl"
)

type Module interface {
	RegisterHTTP(mux *http.ServeMux)
	Client() foundationclient.Client
}

type module struct {
	logger           *slog.Logger
	personController controller.PersonController
	client           foundationclient.Client
}

func New(logger *slog.Logger, pool *pgxpool.Pool) Module {
	personRepository := repositoryimpl.NewPersonRepository(pool)
	personMapper := mapperimpl.NewPersonMapper()
	personValidator := validatorimpl.NewPersonValidator()
	personService := serviceimpl.NewPersonService(personRepository, personMapper, personValidator)

	return &module{
		logger:           logger,
		personController: controllerimpl.NewPersonController(personService),
		client:           clientimpl.NewPersonClient(personService),
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

func (service *module) Client() foundationclient.Client {
	return service.client
}
