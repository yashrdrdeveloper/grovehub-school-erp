package projectpool

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/stdlib"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/postgres/operations/projectpool/constants"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/postgres/operations/projectpool/helpers"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/postgres/operations/projectpool/validators"
)

// OpenProjectPool verifies the explicit database and role before returning a pool.
// Errors intentionally omit the supplied connection URI and driver details.
func OpenProjectPool(ctx context.Context, dsn string) (*sql.DB, error) {
	config, err := validators.ProjectConfig(dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, constants.ConnectionTimeout)
	defer cancel()

	db := stdlib.OpenDB(*config)
	helpers.ConfigurePool(db)

	if !helpers.MatchesIdentity(ctx, db, config.Database, config.User) {
		_ = db.Close()
		return nil, errors.New("project PostgreSQL connection or identity verification failed")
	}

	return db, nil
}
