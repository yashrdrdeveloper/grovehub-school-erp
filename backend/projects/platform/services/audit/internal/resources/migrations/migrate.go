package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"grovehub/services/audit/internal/resources"
)

func Migrate(ctx context.Context, db *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil { return err }
	goose.SetBaseFS(resources.Migrations)
	defer goose.SetBaseFS(nil)
	goose.SetTableName("audit_schema_state")
	return goose.UpContext(ctx, db, "migrations")
}
