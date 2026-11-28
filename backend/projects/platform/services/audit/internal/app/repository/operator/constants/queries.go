package constants

const (
	InsertActivity    = "INSERT INTO audit.operator_activities(event_id,source,actor_id,resource_id,operation,occurred_at,correlation_id) VALUES(?,'internal-tool',?,?,?,?,?) ON CONFLICT(event_id) DO NOTHING"
	IdenticalActivity = "SELECT EXISTS(SELECT 1 FROM audit.operator_activities WHERE event_id=? AND source='internal-tool' AND actor_id=? AND resource_id=? AND operation=? AND occurred_at=? AND correlation_id IS NOT DISTINCT FROM CAST(? AS uuid))"
	InsertReceipt     = "INSERT INTO audit.school_mutation_receipts(execution_id,school_id,previous_version,correlation_id,occurred_at,operation) VALUES(?,?,?,?,?,?) ON CONFLICT(execution_id) DO NOTHING"
	IdenticalReceipt  = "SELECT EXISTS(SELECT 1 FROM audit.school_mutation_receipts WHERE execution_id=? AND school_id=? AND previous_version=? AND correlation_id IS NOT DISTINCT FROM CAST(? AS uuid) AND occurred_at=? AND operation=?)"
	SchoolExists      = "SELECT EXISTS(SELECT 1 FROM audit.activities WHERE tenant_id=?)"
	SearchActivities  = "SELECT event_id,actor_id,resource_id,operation,occurred_at,correlation_id FROM audit.operator_activities ORDER BY occurred_at DESC,event_id DESC LIMIT ? OFFSET ?"
	SearchReceipts    = "SELECT execution_id,school_id,previous_version,correlation_id,occurred_at,operation FROM audit.school_mutation_receipts WHERE TRUE"
	SchoolFilter      = " AND school_id = ?"
	OperationFilter   = " AND operation = ?"
	ReceiptOrder      = " ORDER BY occurred_at DESC,execution_id DESC LIMIT ? OFFSET ?"
)
