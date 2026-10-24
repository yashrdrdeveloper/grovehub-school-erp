package mailprovider

import (
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	coretelemetry "grovehub/libs/core/telemetry"
	mailproviderclient "grovehub/services/mailProvider/client"
	clientimpl "grovehub/services/mailProvider/client/impl"
	"grovehub/services/mailProvider/controller"
	controllerimpl "grovehub/services/mailProvider/controller/impl"
	mapperimpl "grovehub/services/mailProvider/mappers/impl"
	repositoryimpl "grovehub/services/mailProvider/repository/impl"
	"grovehub/services/mailProvider/routes"
	serviceimpl "grovehub/services/mailProvider/services/impl"
	transportimpl "grovehub/services/mailProvider/transport/impl"
	validatorimpl "grovehub/services/mailProvider/validators/impl"
)

type Module interface {
	RegisterHTTP(mux *http.ServeMux)
	Client() mailproviderclient.Client
}

type module struct {
	logger            *slog.Logger
	domainController  controller.DomainController
	senderController  controller.SenderController
	messageController controller.MessageController
	client            mailproviderclient.Client
}

func New(logger *slog.Logger, pool *pgxpool.Pool) Module {
	domainRepository := repositoryimpl.NewDomainRepository(pool)
	domainValidator := validatorimpl.NewDomainValidator()
	domainService := serviceimpl.NewDomainService(domainRepository, domainValidator)
	domainMapper := mapperimpl.NewDomainMapper()

	senderRepository := repositoryimpl.NewSenderRepository(pool)
	senderValidator := validatorimpl.NewSenderValidator()
	senderService := serviceimpl.NewSenderService(senderRepository, domainService, senderValidator)
	senderMapper := mapperimpl.NewSenderMapper()

	messageRepository := repositoryimpl.NewMessageRepository(pool)
	messageValidator := validatorimpl.NewMessageValidator()
	messageTransport := transportimpl.NewHarakaSenderFromEnvironment()
	messageService := serviceimpl.NewMessageService(messageRepository, domainService, senderService, messageTransport, messageValidator)
	messageMapper := mapperimpl.NewMessageMapper()

	return &module{
		logger:            logger,
		domainController:  controllerimpl.NewDomainController(domainService, domainMapper),
		senderController:  controllerimpl.NewSenderController(senderService, senderMapper),
		messageController: controllerimpl.NewMessageController(messageService, messageMapper),
		client:            clientimpl.NewClient(messageService, messageMapper),
	}
}

func (service *module) RegisterHTTP(mux *http.ServeMux) {
	routes.RegisterHTTP(mux, service.domainController, service.senderController, service.messageController)
	service.logger.Info(
		"mail provider service registered",
		"service", coretelemetry.ServiceMailProvider.Name,
		"serviceKey", coretelemetry.ServiceMailProvider.Key,
	)
}

func (service *module) Client() mailproviderclient.Client { return service.client }
