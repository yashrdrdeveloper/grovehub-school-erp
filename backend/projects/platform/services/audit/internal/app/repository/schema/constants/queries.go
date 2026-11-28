package constants

const (
	ActivityQuery      = "SELECT subject_account_id, occurrence_id FROM audit.activities LIMIT 0"
	RejectionQuery     = "SELECT fingerprint FROM audit.rejected_events LIMIT 0"
	ReceiptQuery       = "SELECT source_event_id FROM audit.event_receipts LIMIT 0"
	OperatorQuery      = "SELECT event_id,correlation_id FROM audit.operator_activities LIMIT 0"
	SchoolReceiptQuery = "SELECT execution_id,operation FROM audit.school_mutation_receipts LIMIT 0"
)
