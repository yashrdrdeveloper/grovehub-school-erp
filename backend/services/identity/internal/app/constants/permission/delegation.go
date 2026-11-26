package permission

// SchoolDelegable explicitly allows school-owned capabilities. Registry and
// entitlement administration are reserved for a separately authorized workload.
// Unknown and newly introduced capabilities are denied until classified.
func SchoolDelegable(code Code) bool {
	switch code {
	case AccountCreate, AccountSearch, AccountView, AccountUpdate,
		AccountSuspend, AccountActivate, AccountDeactivate, AccountResendInvitation,
		AccountSelfView, RoleView, RoleCreate, RoleUpdate, RoleDeactivate,
		RolePermissionView, RolePermissionManage, AccountRoleView, AccountRoleAssign,
		PermissionView, FeatureView, AuditView:
		return true
	default:
		return false
	}
}
