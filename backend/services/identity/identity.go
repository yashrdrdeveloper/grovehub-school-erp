package identity

import (
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	coretelemetry "grovehub/libs/core/telemetry"
	foundationclient "grovehub/services/foundation/client"
	"grovehub/services/identity/controller"
	controllerimpl "grovehub/services/identity/controller/impl"
	mapperimpl "grovehub/services/identity/mappers/impl"
	repositoryimpl "grovehub/services/identity/repository/impl"
	"grovehub/services/identity/routes"
	serviceimpl "grovehub/services/identity/services/impl"
	validatorimpl "grovehub/services/identity/validators/impl"
)

type Module interface {
	RegisterHTTP(mux *http.ServeMux)
}

type module struct {
	logger                 *slog.Logger
	registrationController controller.RegistrationController
}

func New(logger *slog.Logger, pool *pgxpool.Pool, foundationClient foundationclient.Client) Module {
	accountRepository := repositoryimpl.NewAccountRepository(pool)
	accountValidator := validatorimpl.NewAccountValidator()
	accountService := serviceimpl.NewAccountService(accountRepository, accountValidator)

	passwordCredentialRepository := repositoryimpl.NewPasswordCredentialRepository(pool)
	passwordCredentialValidator := validatorimpl.NewPasswordCredentialValidator()
	passwordCredentialService := serviceimpl.NewPasswordCredentialService(
		passwordCredentialRepository,
		passwordCredentialValidator,
	)

	registrationMapper := mapperimpl.NewRegistrationMapper()
	registrationValidator := validatorimpl.NewRegistrationValidator()
	registrationService := serviceimpl.NewRegistrationService(
		foundationClient,
		accountService,
		passwordCredentialService,
		registrationMapper,
		registrationValidator,
	)

	return &module{
		logger:                 logger,
		registrationController: controllerimpl.NewRegistrationController(registrationService),
	}
}

func (service *module) RegisterHTTP(mux *http.ServeMux) {
	routes.RegisterHTTP(mux, service.registrationController)
	service.logger.Info(
		"identity service registered",
		"service", coretelemetry.ServiceIdentity.Name,
		"serviceKey", coretelemetry.ServiceIdentity.Key,
	)
}
