// Package infrastructure exposes the service bootstrap and compatibility entry points.
package infrastructure

import (
	"context"
	"log/slog"
	"gorm.io/gorm"
	"grovehub/services/infrastructure/client/messaging"
	"grovehub/services/infrastructure/internal/app/bootstrap"
	bootstrapimpl "grovehub/services/infrastructure/internal/app/bootstrap/impl"

)

type Module = bootstrap.Runtime
type Dependencies = bootstrap.Dependencies
type Service = bootstrap.Service

// Bootstrap is the common typed service entry point.
func Bootstrap() Service {
	return bootstrapimpl.NewServiceBootstrap()
}

// New is a compatibility adapter for callers migrating to Bootstrap().Build.
func New(config messaging.Config, db *gorm.DB, logger *slog.Logger) (Module, error) {
	return Bootstrap().Build(Dependencies{Config: config, DB: db, Logger: logger})
}

// Migrate is a compatibility adapter for callers migrating to Bootstrap().Migration.
func Migrate(ctx context.Context, db *gorm.DB) error {
	return Bootstrap().Migration().Run(ctx, db)
}

func NewEventStore(db *gorm.DB) messaging.EventStore {
	return bootstrapimpl.NewEventStore(db)
}

func NewReplaySource(db *gorm.DB) messaging.ReplaySource {
	return bootstrapimpl.NewReplaySource(db)
}
