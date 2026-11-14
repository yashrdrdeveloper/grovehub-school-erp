// Package infrastructure exposes Infrastructure's lifecycle and shared
// messaging contracts. All persistence, dispatch and provider code is private.
package infrastructure

import (
    "context"
    "errors"
    "fmt"
    "log/slog"

    "gorm.io/gorm"

    messagingclient "grovehub/services/infrastructure/client/messaging"
    bootstrap "grovehub/services/infrastructure/internal/app/bootstrap"
    bootstrapimpl "grovehub/services/infrastructure/internal/app/bootstrap/impl"
    migrations "grovehub/services/infrastructure/internal/resources/migrations"
)

type Module = bootstrap.Runtime

func New(config messagingclient.Config, db *gorm.DB, logger *slog.Logger) (Module,error) {
    return bootstrapimpl.New(config,db,logger)
}

// Migrate uses Infrastructure's independent, embedded Goose migration history.
func Migrate(ctx context.Context,db *gorm.DB) error {
    if db == nil {return errors.New("infrastructure migrations require a database")}
    sqlDB,err:=db.DB()
    if err!=nil{return fmt.Errorf("access infrastructure migration database: %w",err)}
    return migrations.Migrate(ctx,sqlDB)
}
