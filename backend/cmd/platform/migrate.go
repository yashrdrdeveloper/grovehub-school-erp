package main

import (
    "context"
    "fmt"
    "os"
    "time"

    "grovehub/launcher/bootstrap"
    "grovehub/services/foundation"
)

// runMigrationCommand supports embedded Foundation migrations in production,
// without requiring migration SQL files or an external Goose executable.
func runMigrationCommand() int {
    if len(os.Args) != 3 || os.Args[1] != "migrate" || os.Args[2] != "foundation" {
        _, _ = fmt.Fprintln(os.Stderr, "usage: groovehub-platform migrate foundation")
        return 2
    }
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
    defer cancel()

    db, err := bootstrap.NewPostgresDB(ctx)
    if err != nil {
        _, _ = fmt.Fprintf(os.Stderr, "open migration database: %v\n", err)
        return 1
    }
    sqlDB, err := db.DB()
    if err != nil {
        _, _ = fmt.Fprintf(os.Stderr, "access migration database: %v\n", err)
        return 1
    }
    defer func() { _ = sqlDB.Close() }()
    if err := foundation.Migrate(ctx, db); err != nil {
        _, _ = fmt.Fprintf(os.Stderr, "migrate foundation: %v\n", err)
        return 1
    }
    return 0
}
