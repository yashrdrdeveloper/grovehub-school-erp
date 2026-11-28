package bootstrap

import (
	"context"
	"log/slog"
	"net/http"
	"sync/atomic"

	"gorm.io/gorm"

	"grovehub/launcher/internal/src/server"
	"grovehub/libs/platform/session"
	"grovehub/services/audit"
	"grovehub/services/foundation"
	person "grovehub/services/foundation/client/person"
	"grovehub/services/identity"
	"grovehub/services/infrastructure"
	"grovehub/services/infrastructure/client/messaging"
	"grovehub/services/integration"
	connection "grovehub/services/integration/client/connection"
)

// Dependencies contains the external composition adapters owned by the
// Launcher. It keeps lifecycle orchestration independent of transport setup.
type Dependencies struct {
	OpenDatabase          func(context.Context) (*gorm.DB, error)
	OpenAuditDatabase     func(context.Context, *gorm.DB) (*gorm.DB, func(), error)
	NewInfrastructure     func(messaging.Config, *gorm.DB, *slog.Logger) (infrastructure.Module, error)
	NewFoundation         func(*slog.Logger, *gorm.DB, messaging.EventStore) (foundation.Module, error)
	NewIntegration        func(*slog.Logger, *gorm.DB, messaging.MessagingClient, session.TenantStateLookup) (integration.Module, error)
	NewIdentity           func(*slog.Logger, *gorm.DB, person.Client, connection.Client, messaging.EventStore, ...session.TenantStateLookup) (identity.Module, error)
	NewAudit              func(*slog.Logger, *gorm.DB, messaging.MessagingClient, session.ProtectedRouteAuthorizer) (audit.Module, error)
	StartTenantReconciler func(context.Context, *slog.Logger, *gorm.DB, session.TenantStateLookup) func()
	NewServer             func(string, http.Handler, *slog.Logger, *atomic.Bool) server.Server
}
