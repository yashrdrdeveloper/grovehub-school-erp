package impl

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	platformbootstrap "grovehub/libs/platform/bootstrap"
	"grovehub/services/infrastructure/internal/app/bootstrap"
	"grovehub/services/infrastructure/internal/resources/migrations"
	"grovehub/services/infrastructure/client/messaging"
	repositoryimpl "grovehub/services/infrastructure/internal/app/repository/messaging/impl"
	mapperimpl "grovehub/services/infrastructure/internal/app/mappers/messaging/impl"
	validatorimpl "grovehub/services/infrastructure/internal/app/validators/messaging/impl"
)

// NewServiceBootstrap joins runtime, command and migration registration.
func NewServiceBootstrap() bootstrap.Service {
	return platformbootstrap.RuntimeRegistration[*gorm.DB, bootstrap.Dependencies, bootstrap.Runtime]{
		Service: NewOperatorBootstrap(runMigration),
		BuildRuntime: func(input bootstrap.Dependencies) (bootstrap.Runtime, error) {
			return New(input.Config, input.DB, input.Logger).Build()
		},
	}
}

func runMigration(ctx context.Context, db *gorm.DB) error {
	if db == nil {
		return errors.New("infrastructure migrations require a database")
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("access infrastructure migration database: %w", err)
	}
	return migrations.Migrate(ctx, sqlDB)
}

func NewEventStore(db *gorm.DB) messaging.EventStore {
	return repositoryimpl.New(db, validatorimpl.New(), mapperimpl.New())
}

func NewReplaySource(db *gorm.DB) messaging.ReplaySource {
	return repositoryimpl.NewReplaySource(db)
}
