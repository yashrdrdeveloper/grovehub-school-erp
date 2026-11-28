package helpers

const (
	EnvironmentCreated = "ENVIRONMENT_CREATED"
	TemplateCreated = "TEMPLATE_CREATED"
	RequestSubmitted = "REQUEST_SUBMITTED"
	RequestAutoApproved = "REQUEST_AUTO_APPROVED"
	RequestApproved = "REQUEST_APPROVED"
	RequestRejected = "REQUEST_REJECTED"
	RequestCancelled = "REQUEST_CANCELLED"
	ExecutionStarted = "EXECUTION_STARTED"
	ExecutionSucceeded = "EXECUTION_SUCCEEDED"
	ExecutionFailed = "EXECUTION_FAILED"
	ExecutionUnknown = "EXECUTION_UNKNOWN"
	ExecutionAbandoned = "EXECUTION_ABANDONED"
	ExecutionReconciled = "EXECUTION_RECONCILED"
	OperatorCreated = "OPERATOR_CREATED"
	OperatorActive = "OPERATOR_ACTIVE"
	OperatorDisabled = "OPERATOR_DISABLED"
	AuditDeliveryRequeued = "AUDIT_DELIVERY_REQUEUED"
	OperatorBootstrapped = "OPERATOR_BOOTSTRAPPED"
	AuthLoginSucceeded = "AUTH_LOGIN_SUCCEEDED"
	AuthLogoutSucceeded = "AUTH_LOGOUT_SUCCEEDED"
)

func ValidOperation(operation string) bool {
	switch operation {
	case EnvironmentCreated, TemplateCreated, RequestSubmitted, RequestAutoApproved, RequestApproved, RequestRejected, RequestCancelled, ExecutionStarted, ExecutionSucceeded, ExecutionFailed, ExecutionUnknown, ExecutionAbandoned, ExecutionReconciled, OperatorCreated, OperatorActive, OperatorDisabled, AuditDeliveryRequeued, OperatorBootstrapped, AuthLoginSucceeded, AuthLogoutSucceeded:
		return true
	}

	return false
}
