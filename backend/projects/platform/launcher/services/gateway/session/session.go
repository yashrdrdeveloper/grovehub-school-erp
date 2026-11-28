// Package session exposes the gateway security contracts and context operations.
// It deliberately does not import the gateway runtime or a business service.
package session

import (
	"context"
	"grovehub/launcher/services/gateway/constants"
	"grovehub/launcher/services/gateway/contracts"
	"grovehub/launcher/services/gateway/operations/principal"
	"grovehub/launcher/services/gateway/operations/tenant"
)

type Principal = contracts.Principal
type AccessVerifier = contracts.AccessVerifier
type ProtectedRouteAuthorizer = contracts.ProtectedRouteAuthorizer
type TenantCheck = contracts.TenantCheck
type TenantState = contracts.TenantState
type TenantStateLookup = contracts.TenantStateLookup

const (
	AccessCookie = constants.AccessCookie
	RefreshCookie = constants.RefreshCookie
	CSRFName = constants.CSRFName
	ResetCookie = constants.ResetCookie
	SetupCookie = constants.SetupCookie
)

var ErrTenantDenied = contracts.ErrTenantDenied

func WithPrincipal(ctx context.Context, value Principal) context.Context {
	return principal.WithPrincipal(ctx, value)
}

func PrincipalFrom(ctx context.Context) (Principal, bool) {
	return principal.PrincipalFrom(ctx)
}

func ActiveTenantCheck(lookup TenantStateLookup) TenantCheck {
	return tenant.ActiveTenantCheck(lookup)
}

func WithTenantCheck(ctx context.Context, check TenantCheck) context.Context {
	return tenant.WithTenantCheck(ctx, check)
}

func TenantGeneration(ctx context.Context, id string) (int64, error) {
	return tenant.TenantGeneration(ctx, id)
}

func CheckTenant(ctx context.Context, id string) error {
	return tenant.CheckTenant(ctx, id)
}

func CheckTenantSession(ctx context.Context, id string, issued int64) error {
	return tenant.CheckTenantSession(ctx, id, issued)
}

func CheckAccountTenant(ctx context.Context, id *string) error {
	return tenant.CheckAccountTenant(ctx, id)
}

func CheckAccountGeneration(ctx context.Context, id *string, issued int64) error {
	return tenant.CheckAccountGeneration(ctx, id, issued)
}
