package foundation

import (
	"log/slog"
	"net/http"

	"gorm.io/gorm"

	coretelemetry "grovehub/libs/core/telemetry"
	coretransaction "grovehub/libs/core/database/transaction"
	foundationclient "grovehub/services/foundation/client"
	clientimpl "grovehub/services/foundation/client/impl"
	protected "grovehub/services/foundation/controller/protected"
	controllerimpl "grovehub/services/foundation/controller/protected/impl"
	personeventimpl "grovehub/services/foundation/events/persons/impl"
	mapperimpl "grovehub/services/foundation/mappers/impl"
	repositoryimpl "grovehub/services/foundation/repository/impl"
	"grovehub/services/foundation/routes"
	serviceimpl "grovehub/services/foundation/services/impl"
	validatorimpl "grovehub/services/foundation/validators/impl"
	infrastructureclient "grovehub/services/infrastructure/client"
)

type Module interface {
	RegisterHTTP(mux *http.ServeMux)
	Client() foundationclient.Client
}

type module struct {
	logger           *slog.Logger
	personController protected.PersonController
	client           foundationclient.Client
}

func New(
	logger *slog.Logger,
	db *gorm.DB,
	messagingClient infrastructureclient.MessagingClient,
) Module {
	personRepository := repositoryimpl.NewPersonRepository(db)
	contactRepository := repositoryimpl.NewContactRepository(db)
	unitOfWork := coretransaction.NewUnitOfWork(db)
	personMapper := mapperimpl.NewPersonMapper()
	personValidator := validatorimpl.NewPersonValidator()
	personService := serviceimpl.NewPersonService(personRepository, contactRepository, unitOfWork, personMapper, personValidator)
	personService = personeventimpl.NewPersonService(personService, messagingClient, logger)
	personController := controllerimpl.NewPersonController(personService, personValidator)

	// The internal client uses the exact same generated routes and controller
	// pipeline as external HTTP callers, without a loopback TCP hop.
	internalMux := http.NewServeMux()
	routes.RegisterHTTP(internalMux, personController)

	return &module{
		logger:           logger,
		personController: personController,
		client:           clientimpl.NewPersonClient(internalMux),
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
