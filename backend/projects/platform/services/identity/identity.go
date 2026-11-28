// Package identity exposes the module entrypoint; implementation is private.
package identity

import (
	"context"
	schoolgateway "grovehub/launcher/session"
	"errors"
	"fmt"
	"log/slog"

	"gorm.io/gorm"

	foundationclient "grovehub/services/foundation/client/person"
	integrationclient "grovehub/services/integration/client/connection"
	infrastructure "grovehub/services/infrastructure/client/messaging"
	"grovehub/services/identity/internal/app/bootstrap"
	bootstrapimpl "grovehub/services/identity/internal/app/bootstrap/impl"
	migrations "grovehub/services/identity/internal/resources/migrations"
)

type Module = bootstrap.Runtime

func New(logger *slog.Logger, db *gorm.DB, foundation foundationclient.Client, integration integrationclient.Client, events infrastructure.EventStore, lookups ...schoolgateway.TenantStateLookup) (Module, error) {
	return bootstrapimpl.New(logger, db, foundation, integration, events, lookups...).Build()
}

func Migrate(ctx context.Context, db *gorm.DB) error {
	if db == nil {
		return errors.New("identity migration requires a database")
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("identity migration DB: %w", err)
	}
	return migrations.Migrate(ctx, sqlDB)
}
