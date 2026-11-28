package impl

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	platformbootstrap "grovehub/libs/platform/bootstrap"
	"grovehub/services/foundation/internal/app/bootstrap"
	"grovehub/services/foundation/internal/resources/migrations"

)

// NewServiceBootstrap joins runtime, command and migration registration.
func NewServiceBootstrap() bootstrap.Service {
	return platformbootstrap.RuntimeRegistration[*gorm.DB, bootstrap.Dependencies, bootstrap.Runtime]{
		Service: NewOperatorBootstrap(runMigration),
		BuildRuntime: func(input bootstrap.Dependencies) (bootstrap.Runtime, error) {
			return New(input.Logger, input.DB, input.Events).Build()
		},
	}
}

func runMigration(ctx context.Context, db *gorm.DB) error {
	if db == nil {
		return errors.New("foundation migrations require a database")
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("access Foundation migration database: %w", err)
	}
	return migrations.Migrate(ctx, sqlDB)
}
