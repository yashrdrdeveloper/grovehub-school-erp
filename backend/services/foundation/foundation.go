// Package foundation exposes only the module bootstrap entry point.
package foundation

import (
    "context"
    "errors"
    "fmt"
    "log/slog"

    "gorm.io/gorm"

    foundationbootstrap "grovehub/services/foundation/internal/app/bootstrap"
    bootstrapimpl "grovehub/services/foundation/internal/app/bootstrap/impl"
    foundationresources "grovehub/services/foundation/internal/resources"
    infrastructureclient "grovehub/services/infrastructure/client"
)

// Module provides Foundation's shared lifecycle and typed Person client.
type Module = foundationbootstrap.Runtime

func New(logger *slog.Logger, db *gorm.DB, messaging infrastructureclient.MessagingClient) (Module, error) {
    return bootstrapimpl.New(logger, db, messaging).Build()
}

// Migrate applies embedded Foundation SQL migrations using the existing Goose history.
// Run before starting the production HTTP server.
func Migrate(ctx context.Context, db *gorm.DB) error {
    if db == nil { return errors.New("foundation migrations require a database") }
    sqlDB, err := db.DB()
    if err != nil { return fmt.Errorf("access Foundation migration database: %w",err) }
    return foundationresources.Migrate(ctx, sqlDB)
}
