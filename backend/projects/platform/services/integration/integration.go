// Package integration exposes the module bootstrap entry point, like Foundation.
package integration

import (
    "context"
    "errors"
    "fmt"
    "log/slog"

    "gorm.io/gorm"

    integrationconfig "grovehub/services/integration/internal/app/config"
    bootstrap "grovehub/services/integration/internal/app/bootstrap"
    bootstrapimpl "grovehub/services/integration/internal/app/bootstrap/impl"
    migrations "grovehub/services/integration/internal/resources/migrations"
    messaging "grovehub/services/infrastructure/client/messaging"
)

type Module = bootstrap.Runtime

// DevAuth is the minimal public bootstrap option for explicitly enabled,
// loopback-only development authentication; provider configuration stays private.
type DevAuth = integrationconfig.DevAuth

func New(logger *slog.Logger, db *gorm.DB, events messaging.MessagingClient, developmentAuth ...DevAuth) (Module, error) {
    settings:=integrationconfig.DevAuth{}
    if len(developmentAuth)>0 {settings=developmentAuth[0]}
    return bootstrapimpl.New(logger,db,events,settings).Build()
}

func Migrate(ctx context.Context, db *gorm.DB) error {
    if db==nil {return errors.New("integration migrations require a database")}
    sqlDB,err:=db.DB()
    if err!=nil {return fmt.Errorf("access integration migration database: %w",err)}
    return migrations.Migrate(ctx,sqlDB)
}
