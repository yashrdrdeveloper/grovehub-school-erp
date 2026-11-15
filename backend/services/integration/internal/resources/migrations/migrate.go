package migrations

import (
    "context"
    "database/sql"
    "errors"
    "fmt"

    "github.com/pressly/goose/v3"
    resources "grovehub/services/integration/internal/resources"
)

// Migrate applies Integration's own Goose migration history.
func Migrate(ctx context.Context, db *sql.DB) error {
    if db==nil{return errors.New("integration migrations require a database")}
    if err:=goose.SetDialect("postgres");err!=nil{return fmt.Errorf("set integration migration dialect: %w",err)}
    goose.SetBaseFS(resources.Migrations)
    defer goose.SetBaseFS(nil)
    goose.SetTableName("integration_schema_state")
    if err:=goose.UpContext(ctx,db,"migrations");err!=nil{return fmt.Errorf("apply Integration migrations: %w",err)}
    return nil
}
