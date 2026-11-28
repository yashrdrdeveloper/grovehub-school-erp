package validators

import (
	"errors"
	"net/url"
	"strings"

	"github.com/jackc/pgx/v5"
)

var errInvalidURI = errors.New("explicit project PostgreSQL URI with host, database and user is required")

func ProjectConfig(dsn string) (*pgx.ConnConfig, error) {
	uri, err := url.Parse(strings.TrimSpace(dsn))
	if err != nil || !validProjectURI(uri) {
		return nil, errInvalidURI
	}

	if !validQuery(uri.RawQuery) {
		return nil, errInvalidURI
	}

	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, errInvalidURI
	}

	if config.Database != strings.TrimPrefix(uri.Path, "/") || config.User != uri.User.Username() {
		return nil, errInvalidURI
	}

	return config, nil
}

func validProjectURI(uri *url.URL) bool {
	return (uri.Scheme == "postgres" || uri.Scheme == "postgresql") &&
		uri.Hostname() != "" &&
		uri.User != nil &&
		uri.User.Username() != "" &&
		strings.Trim(uri.Path, "/") != "" &&
		uri.Fragment == ""
}

func validQuery(raw string) bool {
	query, err := url.ParseQuery(raw)
	if err != nil {
		return false
	}

	for _, key := range []string{"host", "hostaddr", "port", "database", "dbname", "user", "service", "servicefile"} {
		if query.Has(key) {
			return false
		}
	}

	return true
}
