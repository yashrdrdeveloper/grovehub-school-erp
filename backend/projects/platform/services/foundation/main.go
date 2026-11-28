// Package foundation exposes the service bootstrap and compatibility entry points.
package foundation

import (
	"context"
	"log/slog"
	"gorm.io/gorm"
	"grovehub/services/infrastructure/client/messaging"
	"grovehub/services/foundation/internal/app/bootstrap"
	bootstrapimpl "grovehub/services/foundation/internal/app/bootstrap/impl"

)

type Module = bootstrap.Runtime
type Dependencies = bootstrap.Dependencies
type Service = bootstrap.Service

// Bootstrap is the common typed service entry point.
func Bootstrap() Service {
	return bootstrapimpl.NewServiceBootstrap()
}

// New is a compatibility adapter for callers migrating to Bootstrap().Build.
func New(logger *slog.Logger, db *gorm.DB, events messaging.EventStore) (Module, error) {
	return Bootstrap().Build(Dependencies{Logger: logger, DB: db, Events: events})
}

// Migrate is a compatibility adapter for callers migrating to Bootstrap().Migration.
func Migrate(ctx context.Context, db *gorm.DB) error {
	return Bootstrap().Migration().Run(ctx, db)
}
