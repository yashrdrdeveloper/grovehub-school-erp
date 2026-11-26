package approval

type RequestType string

const (
	RoleAccess RequestType = "ROLE_ACCESS"
	PermissionAccess RequestType = "PERMISSION_ACCESS"
	RoleAssignmentChange RequestType = "ROLE_ASSIGNMENT_CHANGE"
	RolePermissionChange RequestType = "ROLE_PERMISSION_CHANGE"
	FeatureFlagChange RequestType = "FEATURE_FLAG_CHANGE"
	PolicyChange RequestType = "APPROVAL_POLICY_CHANGE"
	RoleLifecycleChange RequestType = "ROLE_LIFECYCLE_CHANGE"
	AccountLifecycleChange RequestType = "ACCOUNT_LIFECYCLE_CHANGE"
)

func Known(value RequestType) bool {
	switch value {
	case RoleAccess, PermissionAccess, RoleAssignmentChange, RolePermissionChange, FeatureFlagChange, PolicyChange, RoleLifecycleChange, AccountLifecycleChange:
		return true
	default:
		return false
	}
}

type Status string

const (
	Pending Status = "PENDING"
	Approved Status = "APPROVED"
	Rejected Status = "REJECTED"
	Cancelled Status = "CANCELLED"
	Expired Status = "EXPIRED"
	Applied Status = "APPLIED"
	Failed Status = "FAILED"
	Stale Status = "STALE"
)

type Decision string

const (
	Approve Decision = "APPROVE"
	Reject Decision = "REJECT"
)

const MaxSteps = 20
