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
	PermissionView          Code = "IDENTITY_PERMISSION_VIEW"
	PermissionRegister      Code = "IDENTITY_PERMISSION_REGISTER"
	FeatureView             Code = "IDENTITY_FEATURE_VIEW"
	FeatureManage           Code = "IDENTITY_FEATURE_MANAGE"
	FeatureRegister         Code = "IDENTITY_FEATURE_REGISTER"
)

type Definition struct {
	Code    Code
	Feature feature.Code
}

func Catalog() []Definition {
	return []Definition{
		{Code: AccountCreate, Feature: feature.AccountManagement},
		{Code: AccountSearch, Feature: feature.AccountManagement},
		{Code: AccountView, Feature: feature.AccountManagement},
		{Code: AccountUpdate, Feature: feature.AccountManagement},
		{Code: AccountSuspend, Feature: feature.AccountManagement},
		{Code: AccountActivate, Feature: feature.AccountManagement},
		{Code: AccountDeactivate, Feature: feature.AccountManagement},
		{Code: AccountResendInvitation, Feature: feature.AccountManagement},
		{Code: AccountSelfView},
		{Code: RoleView, Feature: feature.RoleManagement},
		{Code: RoleCreate, Feature: feature.RoleManagement},
		{Code: RoleUpdate, Feature: feature.RoleManagement},
		{Code: RoleDeactivate, Feature: feature.RoleManagement},
		{Code: RolePermissionView, Feature: feature.RoleManagement},
		{Code: RolePermissionManage, Feature: feature.RoleManagement},
		{Code: AccountRoleView, Feature: feature.RoleManagement},
		{Code: AccountRoleAssign, Feature: feature.RoleManagement},
		{Code: PermissionView, Feature: feature.RoleManagement},
		{Code: PermissionRegister},
		{Code: FeatureView},
		{Code: FeatureManage},
		{Code: FeatureRegister},
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
