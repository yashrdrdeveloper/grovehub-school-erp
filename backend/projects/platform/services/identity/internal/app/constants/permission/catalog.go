// Package permission owns capability codes. A known code is not a registered
// or granted permission. Registration and enforcement are separate operations.
package permission

import "grovehub/services/identity/internal/app/constants/feature"

type Code string

const (
	AccountCreate           Code = "IDENTITY_ACCOUNT_CREATE"
	AccountSearch           Code = "IDENTITY_ACCOUNT_SEARCH"
	AccountView             Code = "IDENTITY_ACCOUNT_VIEW"
	AccountUpdate           Code = "IDENTITY_ACCOUNT_UPDATE"
	AccountSuspend          Code = "IDENTITY_ACCOUNT_SUSPEND"
	AccountActivate         Code = "IDENTITY_ACCOUNT_ACTIVATE"
	AccountDeactivate       Code = "IDENTITY_ACCOUNT_DEACTIVATE"
	AccountResendInvitation Code = "IDENTITY_ACCOUNT_RESEND_INVITATION"
	AccountSelfView         Code = "IDENTITY_ACCOUNT_SELF_VIEW"
	RoleView                Code = "IDENTITY_ROLE_VIEW"
	RoleCreate              Code = "IDENTITY_ROLE_CREATE"
	RoleUpdate              Code = "IDENTITY_ROLE_UPDATE"
	RoleDeactivate          Code = "IDENTITY_ROLE_DEACTIVATE"
	RolePermissionView      Code = "IDENTITY_ROLE_PERMISSION_VIEW"
	RolePermissionManage    Code = "IDENTITY_ROLE_PERMISSION_MANAGE"
	AccountRoleView         Code = "IDENTITY_ACCOUNT_ROLE_VIEW"
	AccountRoleAssign       Code = "IDENTITY_ACCOUNT_ROLE_ASSIGN"
	AuditView               Code = "IDENTITY_AUDIT_VIEW"
	PermissionView          Code = "IDENTITY_PERMISSION_VIEW"
	PermissionRegister      Code = "IDENTITY_PERMISSION_REGISTER"
	FeatureView             Code = "IDENTITY_FEATURE_VIEW"
	FeatureManage           Code = "IDENTITY_FEATURE_MANAGE"
	FeatureRegister         Code = "IDENTITY_FEATURE_REGISTER"
)

type Definition struct {
	Code    Code
	Name    string
	Feature feature.Code
}

func Catalog() []Definition {
	return []Definition{
		{Code: AccountCreate, Name: "Create accounts", Feature: feature.AccountManagement},
		{Code: AccountSearch, Name: "Search accounts", Feature: feature.AccountManagement},
		{Code: AccountView, Name: "View accounts", Feature: feature.AccountManagement},
		{Code: AccountUpdate, Name: "Update accounts", Feature: feature.AccountManagement},
		{Code: AccountSuspend, Name: "Suspend accounts", Feature: feature.AccountManagement},
		{Code: AccountActivate, Name: "Activate accounts", Feature: feature.AccountManagement},
		{Code: AccountDeactivate, Name: "Deactivate accounts", Feature: feature.AccountManagement},
		{Code: AccountResendInvitation, Name: "Resend invitations", Feature: feature.AccountManagement},
		{Code: AccountSelfView, Name: "View own account"},
		{Code: RoleView, Name: "View roles", Feature: feature.RoleManagement},
		{Code: RoleCreate, Name: "Create roles", Feature: feature.RoleManagement},
		{Code: RoleUpdate, Name: "Update roles", Feature: feature.RoleManagement},
		{Code: RoleDeactivate, Name: "Deactivate roles", Feature: feature.RoleManagement},
		{Code: RolePermissionView, Name: "View role permissions", Feature: feature.RoleManagement},
		{Code: RolePermissionManage, Name: "Manage role permissions", Feature: feature.RoleManagement},
		{Code: AccountRoleView, Name: "View account roles", Feature: feature.RoleManagement},
		{Code: AccountRoleAssign, Name: "Assign account roles", Feature: feature.RoleManagement},
		{Code: AuditView, Name: "View tenant RBAC audit history", Feature: feature.AuditManagement},
		{Code: PermissionView, Name: "View permissions", Feature: feature.RoleManagement},
		{Code: PermissionRegister, Name: "Register permissions"},
		{Code: FeatureView, Name: "View features"},
		{Code: FeatureManage, Name: "Manage features"},
		{Code: FeatureRegister, Name: "Register features"},
	}
}

func Find(code Code) (Definition, bool) {
	for _, item := range Catalog() {
		if item.Code == code {
			return item, true
		}
	}
	return Definition{}, false
}
