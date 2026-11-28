package impl

import (
	"context"
	"errors"

	"gorm.io/gorm"

	platformbootstrap "grovehub/libs/platform/bootstrap"
	"grovehub/services/registry/internal/app/bootstrap/constants"
	"grovehub/services/registry/internal/resources/migrations"
)

func New() platformbootstrap.Service[*gorm.DB] {
	return platformbootstrap.Registration[*gorm.DB]{
		ServiceName: constants.ServiceName,
		MigrationPlan: platformbootstrap.Migration[*gorm.DB]{
			Target: platformbootstrap.BusinessDatabase,
			Run: runMigration,
		},
	}
}

func runMigration(ctx context.Context, db *gorm.DB) error {
	if db == nil {
		return errors.New("registry migrations require a database")
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return migrations.Migrate(ctx, sqlDB)
}
