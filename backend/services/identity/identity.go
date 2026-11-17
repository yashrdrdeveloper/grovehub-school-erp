package identity

import (
	"log/slog"
    "context"
    "encoding/base64"
    "time"
	"net/http"
    "os"

	"gorm.io/gorm"

	coretelemetry "grovehub/libs/core/telemetry"
	foundationclient "grovehub/services/foundation/client/person"
	"grovehub/services/identity/controller"
    publiccontroller "grovehub/services/identity/controller/public"
    publiccontrollerimpl "grovehub/services/identity/controller/public/impl"
    identityservices "grovehub/services/identity/services"
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
    Start(ctx context.Context)
}

type module struct {
	logger         *slog.Logger
	authController controller.AuthController
    publicAuthController publiccontroller.AuthController
    protectedMiddleware func(http.Handler) http.Handler
    invitationWorker *serviceimpl.PasswordSetupMailWorker
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

	sessionAuthenticator := serviceimpl.NewSessionAuthenticator(db)
    authMapper := mapperimpl.NewAuthMapper()
	authValidator := validatorimpl.NewAuthValidator()
	dependencies := identityservices.AuthDependencies{
        RegistrationAuthorizer: serviceimpl.NewAdminRegistrationAuthorizer(os.Getenv("IDENTITY_ADMIN_ACCOUNT_IDS")),
        SessionIssuer: serviceimpl.NewSessionIssuer(db),
        PasswordSetupTokens: repositoryimpl.NewPasswordSetupTokenRepository(db),
    }
    var invitationWorker *serviceimpl.PasswordSetupMailWorker
    key, keyErr := base64.StdEncoding.DecodeString(os.Getenv("IDENTITY_INVITATION_KEY_B64"))
    if keyErr == nil && len(key) == 32 {
        queue, queueErr := serviceimpl.NewPasswordSetupInvitationQueue(db, dependencies.PasswordSetupTokens, key, os.Getenv("IDENTITY_SETUP_URL"))
        sender, senderErr := serviceimpl.NewIntegrationMailSender(os.Getenv("IDENTITY_SYSTEM_MAIL_BASE_URL"), os.Getenv("IDENTITY_SYSTEM_MAIL_CONNECTION_ID"), os.Getenv("IDENTITY_SYSTEM_MAIL_ACCESS_TOKEN"))
        if queueErr == nil && senderErr == nil {
            worker, workerErr := serviceimpl.NewPasswordSetupMailWorker(db, key, sender)
            if workerErr == nil { dependencies.InvitationQueue = queue; invitationWorker = worker }
        }
    }
    authService := serviceimpl.NewAuthService(
		foundationClient,
		accountService,
		passwordCredentialService,
		authMapper,
		authValidator,
        db,
        dependencies,
	)

	return &module{
		logger:         logger,
		authController: controllerimpl.NewAuthController(authService, authValidator),
        publicAuthController: publiccontrollerimpl.NewAuthController(authService),
        protectedMiddleware: sessionAuthenticator.Middleware,
        invitationWorker: invitationWorker,
	}
}

func (service *module) RegisterHTTP(mux *http.ServeMux) {
	routes.RegisterHTTP(mux, service.authController, service.publicAuthController, service.protectedMiddleware)
	service.logger.Info(
		"identity service registered",
		"service", coretelemetry.ServiceIdentity.Name,
		"serviceKey", coretelemetry.ServiceIdentity.Key,
	)
}


// Start runs the transactional invitation delivery worker only when a system
// sender and encryption configuration were explicitly supplied. The caller
// controls worker cancellation through the process lifecycle context.
func (service *module) Start(ctx context.Context) {
    if service.invitationWorker == nil {
        service.logger.Warn("identity invitation delivery disabled; protected registration fails closed")
        return
    }
    go func() {
        ticker := time.NewTicker(10 * time.Second)
        defer ticker.Stop()
        for {
            select {
            case <-ctx.Done(): return
            case <-ticker.C:
                child, cancel := context.WithTimeout(ctx, 20*time.Second)
                _, err := service.invitationWorker.RunOnce(child)
                cancel()
                if err != nil { service.logger.Warn("identity invitation delivery attempt failed") }
            }
        }
    }()
}
