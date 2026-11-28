package constants

const (
	Pattern = "POST /audit/restricted/operator-activities"
	SearchPattern = "POST /audit/restricted/operator-activities/search"
	SchoolCheckPattern = "GET /audit/restricted/schools/{schoolId}/activity-check"
	ReceiptPattern = "POST /audit/restricted/school-rename-receipts"
	StatusReceiptPattern = "POST /audit/restricted/school-status-receipts"
	ReceiptSearchPattern = "POST /audit/restricted/school-mutation-receipts/search"

	Scope = "internal.audit.operator.write"
	ReadScope = "internal.audit.operator.read"
	SchoolCheckScope = "internal.audit.school.activity-check"
	ReceiptScope = "internal.audit.school.rename-receipt.write"
	StatusReceiptScope = "internal.audit.school.status-receipt.write"
	ReceiptReadScope = "internal.audit.school.mutation-receipt.read"

	WriterSubject = "internal-tool-audit"
	ReaderSubject = "internal-tool-audit-reader"
	SchoolSubject = "internal-tool-school:"
	ReceiptIssuer = "school-registry-audit"
	StatusReceiptIssuer = "school-registry-status-audit"
	ReceiptReader = "internal-tool-school-receipt-reader"
	SchoolIDParam = "schoolId"
	MaxBodyBytes = 4096
)
