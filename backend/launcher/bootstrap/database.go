package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"gorm.io/gorm"

	corepostgres "grovehub/libs/core/database/postgres"
)

const databaseURLEnv = "DATABASE_URL"

var ErrDatabaseURLRequired = errors.New("DATABASE_URL is required")

func NewPostgresDatabase(ctx context.Context) (*gorm.DB, error) {
	databaseURL := strings.TrimSpace(os.Getenv(databaseURLEnv))
	if databaseURL == "" {
		return nil, ErrDatabaseURLRequired
	}

	database, err := corepostgres.Open(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open postgres database: %w", err)
	}
	return database, nil
}
