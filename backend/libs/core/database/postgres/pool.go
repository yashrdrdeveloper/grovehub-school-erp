package postgres

import (
	"context"
	"database/sql"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

// OpenProjectPool opens an explicitly named database and role. It never uses a
// default DSN and never exposes credentials in returned errors. Provisioning
// and migrations remain the responsibility of the owning application.
func OpenProjectPool(ctx context.Context, dsn string) (*sql.DB, error) {
	config, err := projectConfig(dsn)
	if err != nil { return nil, err }
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	db := stdlib.OpenDB(*config)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(2)
	db.SetConnMaxIdleTime(5*time.Minute)
	db.SetConnMaxLifetime(30*time.Minute)
	var database, role string
	if err := db.QueryRowContext(ctx, "SELECT current_database(), current_user").Scan(&database, &role); err != nil || database != config.Database || role != config.User {
		_ = db.Close()
		return nil, errors.New("project PostgreSQL connection or identity verification failed")
	}
	return db, nil
}

func projectConfig(dsn string) (*pgx.ConnConfig, error) {
	invalid := errors.New("explicit project PostgreSQL URI with host, database and user is required")
	u, err := url.Parse(strings.TrimSpace(dsn))
	if err != nil || (u.Scheme != "postgres" && u.Scheme != "postgresql") || u.Hostname() == "" || u.User == nil || u.User.Username() == "" || strings.Trim(u.Path, "/") == "" || u.Fragment != "" {
		return nil, invalid
	}
	// Do not allow a query parameter to silently replace the target identity.
	q, err := url.ParseQuery(u.RawQuery)
	if err != nil { return nil, invalid }
	for _, key := range []string{"host", "hostaddr", "port", "database", "dbname", "user", "service", "servicefile"} {
		if q.Has(key) { return nil, invalid }
	}
	config, err := pgx.ParseConfig(dsn)
	if err != nil || config.Database != strings.TrimPrefix(u.Path, "/") || config.User != u.User.Username() {
		return nil, invalid
	}
	return config, nil
}
