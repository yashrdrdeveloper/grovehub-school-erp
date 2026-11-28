package impl

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	platformbootstrap "grovehub/libs/platform/bootstrap"

	"grovehub/services/registry/internal/app/bootstrap"
	"grovehub/services/registry/internal/resources/migrations"
)

func NewServiceBootstrap() bootstrap.Service {
	return platformbootstrap.RuntimeRegistration[*gorm.DB, bootstrap.Dependencies, bootstrap.Runtime]{
		Service: NewOperatorBootstrap(runMigration),
		BuildRuntime: Build,
	}
}

func runMigration(ctx context.Context, db *gorm.DB) error {
	if db == nil {
		return errors.New("registry migrations require a database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("access Registry migration database: %w", err)
	}

	return migrations.Migrate(ctx, sqlDB)
}
