package constants

import "errors"

var (
	ErrRequired = errors.New("tenant-id, event-id, database and operator for apply are required")
	ErrAuditDatabaseRequired = errors.New("audit-database is required with AUDIT_DATABASE_URL")
	ErrBusinessConnection = errors.New("could not connect to business database")
	ErrBusinessDatabaseMismatch = errors.New("configured business database does not match --database")
	ErrAuditDatabaseMismatch = errors.New("configured Audit database does not match --audit-database")
	ErrReplay = errors.New("replay failed; verify retained published event, tenant, schema and identity conflicts")
)
