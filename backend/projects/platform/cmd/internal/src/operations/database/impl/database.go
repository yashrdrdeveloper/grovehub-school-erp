package impl

import (
	"context"

	"gorm.io/gorm"

	"grovehub/cmd/internal/src/operations/database"
	"grovehub/launcher"
)

// @Service
type postgres struct{}

func New() database.Database {
	return postgres{}
}

func (postgres) Business(ctx context.Context) (*gorm.DB, func(), error) {
	db, err := launcher.NewPostgresDB(ctx)
	if err != nil {
		return nil, nil, err
	}

	raw, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	return db, func() { _ = raw.Close() }, nil
}

func (postgres) Audit(ctx context.Context, shared *gorm.DB) (*gorm.DB, func(), error) {
	return launcher.NewAuditPostgresDB(ctx, shared)
}
