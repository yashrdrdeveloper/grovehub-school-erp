package impl

import (
	"errors"
	"log/slog"
	"net/url"

	"gorm.io/gorm"
	"grovehub/services/audit"
	auditclient "grovehub/services/audit/client/producer"

	foundationclient "grovehub/services/foundation/client/person"
	"grovehub/services/identity/internal/app/bootstrap"
	configimpl "grovehub/services/identity/internal/app/config/impl"
	securityimpl "grovehub/services/identity/internal/app/services/security/impl"
	integrationclient "grovehub/services/integration/client/connection"
	infrastructure "grovehub/services/infrastructure/client/messaging"
)

type builder struct {
	logger      *slog.Logger
	db          *gorm.DB
	foundation  foundationclient.Client
	integration integrationclient.Client
	events      infrastructure.EventStore
	audit       auditclient.Client
}

var _ bootstrap.Bootstrap = (*builder)(nil)

func New(logger *slog.Logger, db *gorm.DB, foundation foundationclient.Client, integration integrationclient.Client, events infrastructure.EventStore) bootstrap.Bootstrap {
	return &builder{logger: logger, db: db, foundation: foundation, integration: integration, events: events, audit: audit.NewProducer(events)}
}

func (b *builder) Build() (bootstrap.Runtime, error) {
	if b.db == nil || b.foundation == nil || b.integration == nil || b.events == nil {
		return nil, errors.New("identity requires PostgreSQL, Foundation, Integration, and Infrastructure")
	}
	if b.logger == nil {
		b.logger = slog.Default()
	}
	settings, err := configimpl.Load()
	if err != nil {
		return nil, err
	}
	if len(settings.EncryptionKey) != 32 {
		return nil, errors.New("identity invitation encryption key must be configured")
	}
	if settings.ApprovalExecutionTenant != "" {
		link, err := url.Parse(settings.SetupURL)
		if err != nil || link.Scheme != "https" || link.Hostname() == "" || link.User != nil || link.Fragment != "" {
			return nil, errors.New("approval execution requires a trusted HTTPS setup URL")
		}
	}
	if settings.LegacyRBACDevelopment {
		var migrated bool
		if err := b.db.Raw("SELECT EXISTS (SELECT 1 FROM identity.rbac_migration_receipts)").Scan(&migrated).Error; err != nil {
			return nil, err
		}
		if migrated {
			return nil, errors.New("legacy-development cannot run against a migrated database")
		}
	}
	repositories := b.newRepository()
	services := b.newServices(repositories, settings)
	controllers := b.newControllers(services)
	module := &runtime{
		logger:   b.logger,
		routes:   controllers,
		client:   b.newClient(controllers),
		verifier: securityimpl.New(repositories.accounts, repositories.sessions, settings.AdminIDs),
		policy: services.policy,
	}
	module.startWorkers(services.worker)
	return module, nil
}
