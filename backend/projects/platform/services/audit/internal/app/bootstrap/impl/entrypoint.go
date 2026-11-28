package impl

import (
	"context"
	"errors"
	"gorm.io/gorm"
	platformbootstrap "grovehub/libs/platform/bootstrap"
	"grovehub/services/audit/internal/app/bootstrap"
	"grovehub/services/audit/internal/resources/migrations"
	activityrepo "grovehub/services/audit/internal/app/repository/activity/impl"
	"grovehub/services/audit/internal/app/services/replay"
	replayimpl "grovehub/services/audit/internal/app/services/replay/impl"
	"grovehub/services/infrastructure/client/messaging"
)

// NewServiceBootstrap joins runtime, command and migration registration.
func NewServiceBootstrap() bootstrap.Service {
	return platformbootstrap.RuntimeRegistration[*gorm.DB, bootstrap.Dependencies, bootstrap.Runtime]{
		Service: NewOperatorBootstrap(runMigration),
		BuildRuntime: func(input bootstrap.Dependencies) (bootstrap.Runtime, error) {
			return New(input.Logger, input.DB, input.Events, input.Access).Build()
		},
	}
}

func runMigration(ctx context.Context, db *gorm.DB) error {
	if db == nil {
		return errors.New("audit migration requires a database")
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return migrations.Migrate(ctx, sqlDB)
}

// Replay executes the operator-only replay operation.
func Replay(ctx context.Context, db *gorm.DB, source messaging.ReplaySource, input replay.Request) (replay.Result, error) {
	return replayimpl.New(db, source, activityrepo.New(db)).Replay(ctx, input)
}
