package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	corepostgres "grovehub/libs/core/database/postgres"
)

const databaseURLEnv = "DATABASE_URL"

var ErrDatabaseURLRequired = errors.New("DATABASE_URL is required")

func NewPostgresPool(ctx context.Context) (*pgxpool.Pool, error) {
	databaseURL := strings.TrimSpace(os.Getenv(databaseURLEnv))
	if databaseURL == "" {
		return nil, ErrDatabaseURLRequired
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
