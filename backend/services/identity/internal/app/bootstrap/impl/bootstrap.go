package impl

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"gorm.io/gorm"

	coregateway "grovehub/libs/core/gateway"
	coresearch "grovehub/libs/core/search"
	foundationclient "grovehub/services/foundation/client/person"
	accountclient "grovehub/services/identity/client/account"
	"grovehub/services/identity/internal/app/bootstrap"
	configimpl "grovehub/services/identity/internal/app/config/impl"
	accountprotected "grovehub/services/identity/internal/app/controller/account/protected"
	accountpublic "grovehub/services/identity/internal/app/controller/account/public"
	authprotected "grovehub/services/identity/internal/app/controller/authentication/protected"
	authpublic "grovehub/services/identity/internal/app/controller/authentication/public"
	"grovehub/services/identity/internal/app/routes"
	"grovehub/services/identity/internal/app/security/verification"
	accountservice "grovehub/services/identity/internal/app/services/account"
	invitation "grovehub/services/identity/internal/app/services/registration/operations/invitation"
	integrationclient "grovehub/services/integration/client/connection"
	infrastructure "grovehub/services/infrastructure/client/messaging"
)

type builder struct {
	logger *slog.Logger
	db *gorm.DB
	foundation foundationclient.Client
	integration integrationclient.Client
	events infrastructure.EventStore
}

func New(logger *slog.Logger, db *gorm.DB, foundation foundationclient.Client, integration integrationclient.Client, events infrastructure.EventStore) bootstrap.Bootstrap {
	return &builder{
		logger: logger,
		db: db,
		foundation: foundation,
		integration: integration,
		events: events,
	}
}

type runtime struct {
	logger *slog.Logger
	accountPublic accountpublic.AccountController
	accountProtected accountprotected.AccountController
	authPublic authpublic.AuthenticationController
	authProtected authprotected.AuthenticationController
	verifier *verification.Verifier
	accountService accountservice.Service
	worker *invitation.Worker
	stop context.CancelFunc
	done chan struct{}
	once sync.Once
}

var _ bootstrap.Runtime = (*runtime)(nil)

func (b *builder) Build() (bootstrap.Runtime, error) {
	if b.db == nil || b.foundation == nil || b.events == nil {
		return nil, errors.New("identity requires PostgreSQL, Foundation, and Infrastructure")
	}
	if b.logger == nil {
		b.logger = slog.Default()
	}

	config, err := configimpl.Load()
	if err != nil {
		return nil, err
	}

	repositories := b.newRepository()
	services := b.newServices(repositories, config)
	publicAuth, protectedAuth, publicAccount, protectedAccount := b.newControllers(services)
	verifier := verification.New(repositories.accounts, repositories.sessions, config.AdminIDs)

	ctx, stop := context.WithCancel(context.Background())
	module := &runtime{
		logger: b.logger,
		accountPublic: publicAccount,
		accountProtected: protectedAccount,
		authPublic: publicAuth,
		authProtected: protectedAuth,
		verifier: verifier,
		accountService: services.account,
		worker: services.worker,
		stop: stop,
		done: make(chan struct{}),
	}
	go module.runWorker(ctx)

	return module, nil
}

func (module *runtime) RegisterHTTP(mux *http.ServeMux) {
	routes.RegisterHTTP(mux, module.authPublic, module.authProtected, module.accountPublic, module.accountProtected)
}

func (module *runtime) VerifyAccess(ctx context.Context, secret string) (coregateway.Principal, error) {
	return module.verifier.VerifyAccess(ctx, secret)
}

func (module *runtime) AccountClient() accountclient.Client {
	return module
}

func (module *runtime) Get(ctx context.Context, accountID string) (accountclient.Account, error) {
	return module.accountService.Get(ctx, accountID)
}

func (module *runtime) Search(ctx context.Context, query coresearch.Request) (coresearch.Result[accountclient.Account], error) {
	return module.accountService.Search(ctx, query)
}

func (module *runtime) Close() {
	module.once.Do(func() {
		module.stop()
		<-module.done
	})
}

func (module *runtime) runWorker(ctx context.Context) {
	defer close(module.done)

	if module.worker == nil {
		return
	}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			runCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
			_, err := module.worker.RunOnce(runCtx)
			cancel()
			if err != nil {
				module.logger.Warn("Identity invitation delivery failed; retry scheduled")
			}
		}
	}
}
