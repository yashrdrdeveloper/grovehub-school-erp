// Package identity exposes the service bootstrap and compatibility entry points.
package identity

import (
	"context"
	"log/slog"
	"gorm.io/gorm"
	"grovehub/libs/platform/session"
	person "grovehub/services/foundation/client/person"
	"grovehub/services/infrastructure/client/messaging"
	connection "grovehub/services/integration/client/connection"
	"grovehub/services/identity/internal/app/bootstrap"
	bootstrapimpl "grovehub/services/identity/internal/app/bootstrap/impl"
	"grovehub/services/identity/client/invitation"
)

type Module = bootstrap.Runtime
type Dependencies = bootstrap.Dependencies
type Service = bootstrap.Service

// Bootstrap is the common typed service entry point.
func Bootstrap() Service {
	return bootstrapimpl.NewServiceBootstrap()
}

// New is a compatibility adapter for callers migrating to Bootstrap().Build.
func New(logger *slog.Logger, db *gorm.DB, foundation person.Client, integration connection.Client, events messaging.EventStore, lookups ...session.TenantStateLookup) (Module, error) {
	return Bootstrap().Build(Dependencies{Logger: logger, DB: db, Foundation: foundation, Integration: integration, Events: events, TenantState: lookups})
}

// Migrate is a compatibility adapter for callers migrating to Bootstrap().Migration.
func Migrate(ctx context.Context, db *gorm.DB) error {
	return Bootstrap().Migration().Run(ctx, db)
}

func NewInvitationDeliveryReader(db *gorm.DB) invitation.DeliveryReader {
	return bootstrapimpl.NewInvitationDeliveryReader(db)
}

func StartTenantSessionReconciler(parent context.Context, logger *slog.Logger, db *gorm.DB, lookup session.TenantStateLookup) func() {
	return bootstrapimpl.StartTenantSessionReconciler(parent, logger, db, lookup)
}
