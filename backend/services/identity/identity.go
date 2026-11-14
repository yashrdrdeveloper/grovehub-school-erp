package identity

import (
	"log/slog"
	"net/http"

	"gorm.io/gorm"

	coretelemetry "grovehub/libs/core/telemetry"
	foundationclient "grovehub/services/foundation/client/person"
	"grovehub/services/identity/controller"
	controllerimpl "grovehub/services/identity/controller/impl"
	accounteventimpl "grovehub/services/identity/events/accounts/impl"
	mapperimpl "grovehub/services/identity/mappers/impl"
	repositoryimpl "grovehub/services/identity/repository/impl"
	"grovehub/services/identity/routes"
	serviceimpl "grovehub/services/identity/services/impl"
	validatorimpl "grovehub/services/identity/validators/impl"
	infrastructureclient "grovehub/services/infrastructure/client/messaging"
)

type Module interface {
	RegisterHTTP(mux *http.ServeMux)
}

type module struct {
	logger         *slog.Logger
	authController controller.AuthController
}

func New(
	logger *slog.Logger,
	db *gorm.DB,
	foundationClient foundationclient.Client,
	messagingClient infrastructureclient.MessagingClient,
) Module {
	accountRepository := repositoryimpl.NewAccountRepository(db)
	accountValidator := validatorimpl.NewAccountValidator()
	accountService := serviceimpl.NewAccountService(accountRepository, accountValidator)
	accountService = accounteventimpl.NewAccountService(accountService, messagingClient, logger)

	passwordCredentialRepository := repositoryimpl.NewPasswordCredentialRepository(db)
	passwordCredentialValidator := validatorimpl.NewPasswordCredentialValidator()
	passwordCredentialService := serviceimpl.NewPasswordCredentialService(
		passwordCredentialRepository,
		passwordCredentialValidator,
	)

	authMapper := mapperimpl.NewAuthMapper()
	authValidator := validatorimpl.NewAuthValidator()
	authService := serviceimpl.NewAuthService(
		foundationClient,
		accountService,
		passwordCredentialService,
		authMapper,
		authValidator,
        db,
	)

	return &module{
		logger:         logger,
		authController: controllerimpl.NewAuthController(authService, authValidator),
	}
}

func (service *module) RegisterHTTP(mux *http.ServeMux) {
	routes.RegisterHTTP(mux, service.authController)
	service.logger.Info(
		"identity service registered",
		"service", coretelemetry.ServiceIdentity.Name,
		"serviceKey", coretelemetry.ServiceIdentity.Key,
	)
}
