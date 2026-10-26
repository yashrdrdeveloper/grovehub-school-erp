package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, databaseURL)
}

func NewDB(databaseURL string) (*gorm.DB, error) {
	return gorm.Open(gormpostgres.Open(databaseURL), &gorm.Config{})
}
