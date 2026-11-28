package constants

import "errors"

var (
	ErrDatabaseRequired = errors.New("database is required")
	ErrConnection       = errors.New("could not connect to configured Audit database")
	ErrInspection       = errors.New("Audit storage inspection failed; verify database name, schema and permissions")
)
