// Package postgres exposes private PostgreSQL connection operations.
package postgres

import (
	"context"
	"database/sql"

	"gorm.io/gorm"

	"grovehub/libs/core/database/internal/src/postgres/operations/connection"
	"grovehub/libs/core/database/internal/src/postgres/operations/projectpool"
)

func NewDB(databaseURL string) (*gorm.DB, error) {
	return connection.NewDB(databaseURL)
}

func OpenProjectPool(ctx context.Context, dsn string) (*sql.DB, error) {
	return projectpool.OpenProjectPool(ctx, dsn)
}
