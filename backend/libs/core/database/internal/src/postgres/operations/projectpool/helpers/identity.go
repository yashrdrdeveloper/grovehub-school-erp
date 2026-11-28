package helpers

import (
	"context"
	"database/sql"
)

func MatchesIdentity(ctx context.Context, db *sql.DB, expectedDatabase, expectedRole string) bool {
	var database, role string

	err := db.QueryRowContext(ctx, "SELECT current_database(), current_user").Scan(&database, &role)

	return err == nil && database == expectedDatabase && role == expectedRole
}
