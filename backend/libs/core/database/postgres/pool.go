package postgres

import (
	"context"
	"database/sql"

	coredatabase "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"
)

func OpenProjectPool(ctx context.Context, dsn string) (*sql.DB, error) {
	return coredatabase.OpenProjectPool(ctx, dsn)
}
