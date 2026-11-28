package impl

import (
	"log/slog"

	"gorm.io/gorm"
	"grovehub/libs/platform/session"
	person "grovehub/services/foundation/client/person"
	"grovehub/services/infrastructure/client/messaging"
	connection "grovehub/services/integration/client/connection"
	"grovehub/launcher/internal/src/bootstrap"
	serverimpl "grovehub/launcher/internal/src/server/impl"
	"grovehub/services/audit"
	"grovehub/services/foundation"
	"grovehub/services/identity"
	"grovehub/services/infrastructure"
	"grovehub/services/integration"
)

// DefaultDependencies is the production composition root. Each adapter has
// explicit ownership; the builder remains responsible for rollback.
func DefaultDependencies() bootstrap.Dependencies {
	return bootstrap.Dependencies{
		OpenDatabase:          NewPostgresDB,
		OpenAuditDatabase:     NewAuditPostgresDB,
		NewInfrastructure: func(config messaging.Config, db *gorm.DB, logger *slog.Logger) (infrastructure.Module, error) {
			return infrastructure.Bootstrap().Build(infrastructure.Dependencies{Config: config, DB: db, Logger: logger})
		},
		NewFoundation: func(logger *slog.Logger, db *gorm.DB, events messaging.EventStore) (foundation.Module, error) {
			return foundation.Bootstrap().Build(foundation.Dependencies{Logger: logger, DB: db, Events: events})
		},
		NewIntegration: func(logger *slog.Logger, db *gorm.DB, events messaging.MessagingClient, lookup session.TenantStateLookup) (integration.Module, error) {
			return integration.Bootstrap().Build(integration.Dependencies{Logger: logger, DB: db, Events: events, TenantState: []session.TenantStateLookup{lookup}})
		},
		NewIdentity: func(logger *slog.Logger, db *gorm.DB, foundationClient person.Client, integrationClient connection.Client, events messaging.EventStore, lookups ...session.TenantStateLookup) (identity.Module, error) {
			return identity.Bootstrap().Build(identity.Dependencies{Logger: logger, DB: db, Foundation: foundationClient, Integration: integrationClient, Events: events, TenantState: lookups})
		},
		NewAudit: func(logger *slog.Logger, db *gorm.DB, events messaging.MessagingClient, access session.ProtectedRouteAuthorizer) (audit.Module, error) {
			return audit.Bootstrap().Build(audit.Dependencies{Logger: logger, DB: db, Events: events, Access: access})
		},
		StartTenantReconciler: identity.StartTenantSessionReconciler,
		NewServer:             serverimpl.New,
	}
}
