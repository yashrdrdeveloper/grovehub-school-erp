package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/gorm"

	corepostgres "grovehub/libs/core/database/postgres"
)

const databaseURLEnv = "DATABASE_URL"

var ErrDatabaseURLRequired = errors.New("DATABASE_URL is required")

func NewPostgresPool(ctx context.Context) (*pgxpool.Pool, error) {
	databaseURL, err := databaseURL()
	if err != nil {
		return nil, err
	}

	pool, err := corepostgres.NewPool(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return pool, nil
}

func NewPostgresDB(ctx context.Context) (*gorm.DB, error) {
	databaseURL, err := databaseURL()
	if err != nil {
		return nil, err
	}

	db, err := corepostgres.NewDB(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("create postgres gorm db: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("access postgres sql db: %w", err)
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("ping postgres gorm db: %w", err)
	}

	return db, nil
}

func databaseURL() (string, error) {
	databaseURL := strings.TrimSpace(os.Getenv(databaseURLEnv))
	if databaseURL == "" {
		return "", ErrDatabaseURLRequired
	}
	return databaseURL, nil
}
