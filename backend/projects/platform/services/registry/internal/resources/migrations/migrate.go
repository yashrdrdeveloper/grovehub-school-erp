package migrations

import (
	"context"
	"database/sql"
	"errors"

	"github.com/pressly/goose/v3"

	"grovehub/services/registry/internal/resources"
)

func Migrate(ctx context.Context, db *sql.DB) error {
	if db == nil {
		return errors.New("registry migrations require a database")
	}
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	goose.SetBaseFS(resources.Migrations)
	defer goose.SetBaseFS(nil)
	goose.SetTableName("registry_schema_state")
	return goose.UpContext(ctx, db, "migrations")
}
