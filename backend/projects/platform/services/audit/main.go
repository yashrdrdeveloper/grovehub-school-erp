// Package audit exposes the service bootstrap and compatibility entry points.
package audit

import (
	"context"
	"log/slog"
	"gorm.io/gorm"
	"grovehub/libs/platform/session"
	"grovehub/services/infrastructure/client/messaging"
	"grovehub/services/audit/internal/app/bootstrap"
	bootstrapimpl "grovehub/services/audit/internal/app/bootstrap/impl"
	producer "grovehub/services/audit/client/producer"
	"grovehub/services/audit/internal/app/services/replay"
)

type Module = bootstrap.Runtime
type Dependencies = bootstrap.Dependencies
type Service = bootstrap.Service

// Bootstrap is the common typed service entry point.
func Bootstrap() Service {
	return bootstrapimpl.NewServiceBootstrap()
}

// New is a compatibility adapter for callers migrating to Bootstrap().Build.
func New(logger *slog.Logger, db *gorm.DB, events messaging.MessagingClient, access session.ProtectedRouteAuthorizer) (Module, error) {
	return Bootstrap().Build(Dependencies{Logger: logger, DB: db, Events: events, Access: access})
}

// Migrate is a compatibility adapter for callers migrating to Bootstrap().Migration.
func Migrate(ctx context.Context, db *gorm.DB) error {
	return Bootstrap().Migration().Run(ctx, db)
}

// NewProducer constructs the Audit producer without starting its runtime.
func NewProducer(store messaging.EventStore) producer.Client {
	return bootstrapimpl.NewProducer(store)
}

type ReplayRequest = replay.Request
type ReplayResult = replay.Result

func Replay(ctx context.Context, db *gorm.DB, source messaging.ReplaySource, input ReplayRequest) (ReplayResult, error) {
	return bootstrapimpl.Replay(ctx, db, source, input)
}
