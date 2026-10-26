package foundation

import (
	"log/slog"
	"net/http"

	"gorm.io/gorm"

	coretelemetry "grovehub/libs/core/telemetry"
	foundationclient "grovehub/services/foundation/client"
	clientimpl "grovehub/services/foundation/client/impl"
	"grovehub/services/foundation/controller"
	controllerimpl "grovehub/services/foundation/controller/impl"
	repositoryimpl "grovehub/services/foundation/repository/impl"
	"grovehub/services/foundation/routes"
	"grovehub/services/foundation/services"
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

func New(logger *slog.Logger, database *gorm.DB) Module {
	personRepository := repositoryimpl.NewPersonRepository(database)
	personService := services.NewPersonService(personRepository)

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
