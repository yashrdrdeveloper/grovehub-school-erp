// Package accounting exposes the module entrypoint; implementation is private.
package accounting

import (
 "context"
 "errors"
 "log/slog"

 "gorm.io/gorm"
 "grovehub/services/accounting/internal/app/bootstrap"
 bootstrapimpl "grovehub/services/accounting/internal/app/bootstrap/impl"
 "grovehub/services/accounting/internal/resources/migrations"
 "grovehub/services/infrastructure/client/messaging"
)

type Module = bootstrap.Runtime

func New(logger *slog.Logger, db *gorm.DB, events messaging.EventStore) (Module, error) {
 return bootstrapimpl.New(logger, db, events).Build()
}

func Migrate(ctx context.Context, db *gorm.DB) error {
 if db == nil { return errors.New("accounting migration requires a database") }
 sqlDB, err := db.DB()
 if err != nil { return err }
 return migrations.Migrate(ctx, sqlDB)
}
