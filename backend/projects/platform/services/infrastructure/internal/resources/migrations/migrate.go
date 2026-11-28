package migrations

import (
    "context"
    "database/sql"
    "errors"
    "fmt"

    "github.com/pressly/goose/v3"

    resources "grovehub/services/infrastructure/internal/resources"
)

// Migrate applies the Infrastructure schema using its independent Goose history.
// Call during deployment before starting any module that writes outbox entries.
func Migrate(ctx context.Context, db *sql.DB) error {
    if db == nil { return errors.New("infrastructure migrations require a database") }
    if err := goose.SetDialect("postgres"); err != nil { return fmt.Errorf("set infrastructure dialect: %w", err) }
    goose.SetBaseFS(resources.Migrations)
    defer goose.SetBaseFS(nil)
    goose.SetTableName("infrastructure_schema_state")
    if err := goose.UpContext(ctx, db, "migrations"); err != nil {
        return fmt.Errorf("apply infrastructure migrations: %w", err)
    }
    return nil
}
