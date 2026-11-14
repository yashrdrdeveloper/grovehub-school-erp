package migrations

import (
    "context"
    "database/sql"
    "errors"
    "fmt"

    "github.com/pressly/goose/v3"

    foundationresources "grovehub/services/foundation/internal/resources"
)

// Migrate runs Foundation's embedded SQL migrations using the same Goose
// state table as the development migration script. Call during deployment,
// before starting Foundation; Goose's global settings are not concurrent-safe.
func Migrate(ctx context.Context, db *sql.DB) error {
    if db == nil { return errors.New("foundation migrations require a database") }
    if err := goose.SetDialect("postgres"); err != nil { return fmt.Errorf("set Foundation migration dialect: %w", err) }
    goose.SetBaseFS(foundationresources.Migrations)
    defer goose.SetBaseFS(nil)
    goose.SetTableName("foundation_schema_state")
    if err := goose.UpContext(ctx, db, "migrations"); err != nil {
        return fmt.Errorf("apply Foundation migrations: %w", err)
    }
    return nil
}
