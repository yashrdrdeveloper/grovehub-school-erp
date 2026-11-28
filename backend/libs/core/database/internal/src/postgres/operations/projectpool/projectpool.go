package projectpool

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"

	"grovehub/libs/core/database/internal/src/postgres/operations/projectpool/constants"
	"grovehub/libs/core/database/internal/src/postgres/operations/projectpool/helpers"
	"grovehub/libs/core/database/internal/src/postgres/operations/projectpool/validators"
)

// OpenProjectPool verifies the explicit database and role before returning a pool.
// Errors intentionally omit the supplied connection URI and driver details.
func OpenProjectPool(ctx context.Context, dsn string) (*sql.DB, error) {
	return OpenVerifiedPool(ctx, dsn, func(config pgx.ConnConfig) *sql.DB { return stdlib.OpenDB(config) })
}

// OpenVerifiedPool owns validation, timeout, pool sizing and identity verification.
// The opener is the SQL-driver boundary; the production facade always uses pgx.
func OpenVerifiedPool(ctx context.Context, dsn string, open func(pgx.ConnConfig) *sql.DB) (*sql.DB, error) {
	config, err := validators.ProjectConfig(dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, constants.ConnectionTimeout)
	defer cancel()

	db := open(*config)
	helpers.ConfigurePool(db)

	if !helpers.MatchesIdentity(ctx, db, config.Database, config.User) {
		_ = db.Close()
		return nil, errors.New("project PostgreSQL connection or identity verification failed")
	}

	return db, nil
}
