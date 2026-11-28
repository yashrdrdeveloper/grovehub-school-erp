package permission

// AdministratorRecoveryCapabilities is the provisional minimum authority needed
// to maintain role grants and account assignments. It is not a role name or an
// implicit grant. Both capabilities must be effective on the same account.
func AdministratorRecoveryCapabilities() []Code {
	return []Code{RolePermissionManage, AccountRoleAssign}
}
