package audit

import (
	"context"
	"errors"
	"log/slog"

	"gorm.io/gorm"
	schoolgateway "grovehub/libs/platform/session"
	producer "grovehub/services/audit/client/producer"
	"grovehub/services/audit/internal/app/bootstrap"
	bootstrapimpl "grovehub/services/audit/internal/app/bootstrap/impl"
	"grovehub/services/audit/internal/resources/migrations"
	activityrepo "grovehub/services/audit/internal/app/repository/activity/impl"
	"grovehub/services/audit/internal/app/services/replay"
	replayimpl "grovehub/services/audit/internal/app/services/replay/impl"
	"grovehub/services/infrastructure/client/messaging"
)

type Module = bootstrap.Runtime

func New(logger *slog.Logger, db *gorm.DB, events messaging.MessagingClient, access schoolgateway.ProtectedRouteAuthorizer) (Module, error) {
	return bootstrapimpl.New(logger, db, events, access).Build()
}

func Migrate(ctx context.Context, db *gorm.DB) error {
	if db == nil { return errors.New("audit migration requires a database") }
	sqlDB, err := db.DB()
	if err != nil { return err }
	return migrations.Migrate(ctx, sqlDB)
}

// NewProducer exposes the Audit-owned event producer without starting a consumer
// or requiring the Audit database. The caller's outbox transaction is preserved.
func NewProducer(store messaging.EventStore) producer.Client {
	return bootstrapimpl.NewProducer(store)
}

// Replay is an operator-only API, not a registered HTTP route.
type ReplayRequest = replay.Request
type ReplayResult = replay.Result

func Replay(ctx context.Context,db *gorm.DB,source messaging.ReplaySource,input ReplayRequest) (ReplayResult,error) {
	return replayimpl.New(db,source,activityrepo.New(db)).Replay(ctx,input)
}
