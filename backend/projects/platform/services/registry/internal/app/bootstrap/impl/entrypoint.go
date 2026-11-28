package impl

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"grovehub/services/registry/internal/app/bootstrap"
	"grovehub/services/registry/internal/resources/migrations"
)

func NewServiceBootstrap() bootstrap.Service {
	return NewOperatorBootstrap(runMigration)
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
