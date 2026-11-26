// Package routepolicy owns the Identity protected-route permission inventory.
package routepolicy

import (
	"context"

	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/services/authorization"
	"grovehub/services/identity/internal/app/services/security/principal"
)

type Rule struct {
	Permission permission.Code
	SelfService bool
	BlockedUntilCutover bool
}

func Lookup(pattern string) (Rule, bool) {
	rules := map[string]Rule{
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
}

func New(authorizer authorization.Service, legacyDevelopment bool) *Policy {
	return &Policy{authorizer: authorizer, legacyDevelopment: legacyDevelopment}
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
		return nil
	}
	if rule.BlockedUntilCutover {
		return authorization.ErrForbidden
	}
	if policy.authorizer == nil {
		return authorization.ErrUnavailable
	}
	return policy.authorizer.Require(ctx, rule.Permission)
}
