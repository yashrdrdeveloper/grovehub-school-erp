package infrastructure

import (
	"log/slog"
	"net/http"

	"gorm.io/gorm"

	infrastructureclient "grovehub/services/infrastructure/client"
	clientimpl "grovehub/services/infrastructure/client/impl"
	infrastructureconfig "grovehub/services/infrastructure/config"
	"grovehub/services/infrastructure/controller"
	controllerimpl "grovehub/services/infrastructure/controller/impl"
	repositoryimpl "grovehub/services/infrastructure/repository/impl"
	"grovehub/services/infrastructure/routes"
	serviceimpl "grovehub/services/infrastructure/services/impl"
)

type Module interface {
	RegisterHTTP(mux *http.ServeMux)
	Client() infrastructureclient.MailClient
}

type module struct {
	logger         *slog.Logger
	dkimController controller.DKIMController
	client         infrastructureclient.MailClient
}

func New(logger *slog.Logger, db *gorm.DB) Module {
	dkimRepository := repositoryimpl.NewDKIMRepository(db)
	dkimService := serviceimpl.NewDKIMService(dkimRepository, infrastructureconfig.DKIMFromEnvironment())
	smtpService := serviceimpl.NewSMTPService(infrastructureconfig.SMTPFromEnvironment())

	return &module{
		logger:         logger,
		dkimController: controllerimpl.NewDKIMController(dkimService),
		client:         clientimpl.NewMailClient(smtpService),
	}
}

func (service *module) RegisterHTTP(mux *http.ServeMux) {
	routes.RegisterHTTP(mux, service.dkimController)
	service.logger.Info("infrastructure service registered")
}

func (service *module) Client() infrastructureclient.MailClient { return service.client }
