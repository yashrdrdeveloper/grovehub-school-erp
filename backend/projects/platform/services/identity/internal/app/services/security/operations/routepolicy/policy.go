// Package routepolicy owns the Identity protected-route permission inventory.
package routepolicy

import (
	"context"

	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/services/authorization"
	"grovehub/services/identity/internal/app/services/approval/operations/cutover"
	"grovehub/services/identity/internal/app/services/security/principal"
)

type Rule struct {
	Permission permission.Code
	SelfService bool
	BlockedUntilCutover bool
	ScopedRead bool
	ScopedMutation bool
	StrictOnly bool
}

func Lookup(pattern string) (Rule, bool) {
	rules := map[string]Rule{
		"POST /audit/protected/search": {ScopedRead: true, StrictOnly: true},
		"AUDIT tenant search": {Permission: permission.AuditView, StrictOnly: true},
		"POST /identity/protected/roles": {Permission: permission.RoleCreate, StrictOnly: true},
		"PATCH /identity/protected/roles/{roleId}": {Permission: permission.RoleUpdate, StrictOnly: true},
		"GET /identity/protected/capabilities/me": {ScopedRead: true, StrictOnly: true},
		"GET /identity/protected/capabilities/permissions": {Permission: permission.PermissionView, StrictOnly: true},
		"POST /identity/protected/approvals/lifecycle": {ScopedMutation: true, StrictOnly: true},
		"POST /identity/protected/approvals/permissions": {Permission: permission.RolePermissionManage, StrictOnly: true},
		"POST /identity/protected/approvals/assignments": {Permission: permission.AccountRoleAssign, StrictOnly: true},
		"POST /identity/protected/approvals/{requestId}/approve": {ScopedMutation: true, StrictOnly: true},
		"POST /identity/protected/approvals/{requestId}/reject": {ScopedMutation: true, StrictOnly: true},
		"POST /identity/protected/approvals/{requestId}/cancel": {ScopedMutation: true, StrictOnly: true},
		"POST /identity/protected/approvals/{requestId}/execute": {BlockedUntilCutover: true, StrictOnly: true},
		"POST /identity/protected/roles/search": {Permission: permission.RoleView, StrictOnly: true},
		"GET /identity/protected/roles/{roleId}": {Permission: permission.RoleView, StrictOnly: true},
		"GET /identity/protected/roles/{roleId}/permissions": {Permission: permission.RolePermissionView, StrictOnly: true},
		"POST /identity/protected/roles/{roleId}/permissions/preview": {Permission: permission.RolePermissionManage, StrictOnly: true},
		"GET /identity/protected/account-roles/{accountId}": {Permission: permission.AccountRoleView, StrictOnly: true},
		"POST /identity/protected/account-roles/{accountId}/preview": {Permission: permission.AccountRoleAssign, StrictOnly: true},
		"POST /identity/protected/approvals/history/search": {ScopedRead: true, StrictOnly: true},
		"POST /identity/protected/approvals/inbox/search": {ScopedRead: true, StrictOnly: true},
		"GET /identity/protected/approvals/inbox/count": {ScopedRead: true, StrictOnly: true},
		"GET /identity/protected/approvals/{requestId}": {ScopedRead: true, StrictOnly: true},
		"POST /identity/protected/accounts": {Permission: permission.AccountCreate, BlockedUntilCutover: true},
		"POST /identity/protected/accounts/search": {Permission: permission.AccountSearch},
		"GET /identity/protected/accounts/roles": {Permission: permission.RoleView},
		"GET /identity/protected/accounts/{accountId}": {Permission: permission.AccountView},
		"PATCH /identity/protected/accounts/{accountId}": {Permission: permission.AccountUpdate, BlockedUntilCutover: true},
		"PATCH /identity/protected/accounts/{accountId}/suspend": {Permission: permission.AccountSuspend, BlockedUntilCutover: true},
		"PATCH /identity/protected/accounts/{accountId}/activate": {Permission: permission.AccountActivate, BlockedUntilCutover: true},
		"PATCH /identity/protected/accounts/{accountId}/deactivate": {Permission: permission.AccountDeactivate, BlockedUntilCutover: true},
		"POST /identity/protected/accounts/{accountId}/resend-invitation": {Permission: permission.AccountResendInvitation, BlockedUntilCutover: true},
		"POST /identity/protected/auth/logout": {SelfService: true},
		"POST /identity/protected/auth/logout-all": {SelfService: true},
		"POST /identity/protected/auth/change-password": {SelfService: true},
		"POST /identity/protected/auth/reauthenticate": {SelfService: true},
		"GET /identity/protected/auth/sessions": {SelfService: true},
		"DELETE /identity/protected/auth/sessions/{sessionId}": {SelfService: true},
		"GET /identity/protected/auth/me": {SelfService: true},
	}
	rule, known := rules[pattern]
	return rule, known
}

type Policy struct {
	authorizer authorization.Service
	legacyDevelopment bool
	execution *cutover.Gate
}

func New(authorizer authorization.Service, legacyDevelopment bool) *Policy {
	return NewWithExecution(authorizer, legacyDevelopment, nil)
}

func NewWithExecution(authorizer authorization.Service, legacyDevelopment bool, gate *cutover.Gate) *Policy {
	return &Policy{authorizer: authorizer, legacyDevelopment: legacyDevelopment, execution: gate}
}

func (policy *Policy) Authorize(ctx context.Context, pattern string) error {
	rule, known := Lookup(pattern)
	if !known {
		return authorization.ErrForbidden
	}
	if rule.SelfService {
		_, err := principal.Require(ctx)
		return err
	}
	if _, err := principal.RequireTenant(ctx); err != nil {
		return err
	}
	if policy.legacyDevelopment {
		if rule.StrictOnly {
			return authorization.ErrForbidden
		}
		return nil
	}
	if rule.BlockedUntilCutover {
		if policy.execution == nil {
			return authorization.ErrForbidden
		}
		switch pattern {
		case "POST /identity/protected/approvals/{requestId}/execute":
			return policy.execution.Require(ctx)
		case "POST /identity/protected/accounts", "POST /identity/protected/accounts/{accountId}/resend-invitation", "PATCH /identity/protected/accounts/{accountId}":
			if err := policy.execution.Require(ctx); err != nil {
				return err
			}
		default:
			return authorization.ErrForbidden
		}
	}
	if rule.ScopedRead || rule.ScopedMutation {
		// The owning approval service enforces live request-specific authority.
		return nil
	}
	if policy.authorizer == nil {
		return authorization.ErrUnavailable
	}
	return policy.authorizer.Require(ctx, rule.Permission)
}
