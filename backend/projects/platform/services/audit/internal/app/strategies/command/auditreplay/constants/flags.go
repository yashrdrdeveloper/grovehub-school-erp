package constants

const (
	TenantFlag = "tenant-id"
	EventFlag = "event-id"
	AuditDatabaseFlag = "audit-database"
	OperatorFlag = "operator"
	ApplyFlag = "apply"

	TenantDescription = "Exact tenant identifier"
	EventDescription = "Exact original published event UUID"
	DatabaseDescription = "Expected business database name"
	AuditDatabaseDescription = "Expected Audit database name; required with AUDIT_DATABASE_URL"
	OperatorDescription = "Operator label (required for apply)"
	ApplyDescription = "Persist the original audit event and replay receipt"
)
