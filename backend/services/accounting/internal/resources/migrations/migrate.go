package migrations

import (
 "context"
 "database/sql"
 "errors"

 "github.com/pressly/goose/v3"
 "grovehub/services/accounting/internal/resources"
)

func Migrate(ctx context.Context, db *sql.DB) error {
 if db == nil { return errors.New("accounting requires a database") }
 if err := goose.SetDialect("postgres"); err != nil { return err }
 goose.SetBaseFS(resources.Migrations)
 defer goose.SetBaseFS(nil)
 goose.SetTableName("accounting_schema_state")
 return goose.UpContext(ctx, db, "migrations")
}
