package foundation

import (
    "context"
    "time"
	"log/slog"
	"net/http"

	"gorm.io/gorm"

	coretelemetry "grovehub/libs/core/telemetry"
	coretransaction "grovehub/libs/core/database/transaction"
    coreoutbox "grovehub/libs/core/events/outbox"
	foundationclient "grovehub/services/foundation/client/person"
	clientimpl "grovehub/services/foundation/client/person/impl"
	protected "grovehub/services/foundation/controller/person/protected"
	controllerimpl "grovehub/services/foundation/controller/person/protected/impl"
    eventpipelines "grovehub/services/foundation/events/person/pipelines"
    eventstrategyimpl "grovehub/services/foundation/events/person/pipelines/strategies/impl"
    eventproducerimpl "grovehub/services/foundation/events/person/producers/impl"
	mapperimpl "grovehub/services/foundation/mappers/person/impl"
	personrepositoryimpl "grovehub/services/foundation/repository/person/impl"
	contactrepositoryimpl "grovehub/services/foundation/repository/contact/impl"
	"grovehub/services/foundation/routes"
	serviceimpl "grovehub/services/foundation/services/person/impl"
	validatorimpl "grovehub/services/foundation/validators/person/impl"
	infrastructureclient "grovehub/services/infrastructure/client"
)

type Module interface {
    RegisterHTTP(mux *http.ServeMux)
    Client() foundationclient.Client
    Close()
}

type module struct {
	logger           *slog.Logger
	personController protected.PersonController
	client           foundationclient.Client
    stopDispatcher context.CancelFunc
    dispatcherDone chan struct{}
}

func New(
	logger *slog.Logger,
	db *gorm.DB,
	messagingClient infrastructureclient.MessagingClient,
) Module {
	personRepository := personrepositoryimpl.NewPersonRepository(db)
	contactRepository := contactrepositoryimpl.NewContactRepository(db)
	unitOfWork := coretransaction.NewUnitOfWork(db)
	personMapper := mapperimpl.NewPersonMapper()
	personValidator := validatorimpl.NewPersonValidator()
    outboxStore := coreoutbox.NewGormStore(db, "foundation.event_outbox")
    eventPipeline := eventpipelines.NewChangedPipeline(eventstrategyimpl.New(), outboxStore)
    eventProducer := eventproducerimpl.New(eventPipeline)
    personService := serviceimpl.NewPersonService(personRepository, contactRepository, unitOfWork, personMapper, personValidator, eventProducer)
    dispatchCtx, stopDispatcher := context.WithCancel(context.Background())
    dispatcherDone := make(chan struct{})
    dispatcher := coreoutbox.NewDispatcher(outboxStore, messagingClient, logger)
    go func() {
        defer close(dispatcherDone)
        dispatcher.Run(dispatchCtx, time.Second)
    }()
	personController := controllerimpl.NewPersonController(personService, personValidator)

	// The internal client uses the exact same generated routes and controller
	// pipeline as external HTTP callers, without a loopback TCP hop.
	internalMux := http.NewServeMux()
	routes.RegisterHTTP(internalMux, personController)

	return &module{
		logger:           logger,
		personController: personController,
		client:           clientimpl.NewPersonClient(internalMux),
        stopDispatcher: stopDispatcher,
        dispatcherDone: dispatcherDone,
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

// Close stops outbox delivery before messaging and database shutdown.
func (service *module) Close() {
    if service.stopDispatcher == nil { return }
    service.stopDispatcher()
    <-service.dispatcherDone
}
