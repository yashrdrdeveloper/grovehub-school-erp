package tenant

import (
	"context"
	"grovehub/internal/session/contracts"
)

func ActiveTenantCheck(lookup contracts.TenantStateLookup) contracts.TenantCheck {
 if lookup == nil { return nil }
 return func(ctx context.Context, tenant string) (int64, error) {
  state, err := lookup(ctx, tenant)
  if err != nil || !state.Active || state.Generation < 1 { return 0, contracts.ErrTenantDenied }
  return state.Generation, nil
 }
}
type tenantCheckKey struct{}

func WithTenantCheck(ctx context.Context, check contracts.TenantCheck) context.Context {
 if check == nil { return ctx }
 return context.WithValue(ctx, tenantCheckKey{}, check)
}

// TenantGeneration is called only with trusted tenant identity. Zero is
// reserved for isolated module/test use where no registry policy is configured.
func TenantGeneration(ctx context.Context, tenant string) (int64, error) {
 check, _ := ctx.Value(tenantCheckKey{}).(contracts.TenantCheck)
 if check == nil { return 0, nil }
 if tenant == "" { return 0, contracts.ErrTenantDenied }
 generation, err := check(ctx, tenant)
 if err != nil || generation < 1 { return 0, contracts.ErrTenantDenied }
 return generation, nil
}

func CheckTenant(ctx context.Context, tenant string) error {
 _, err := TenantGeneration(ctx, tenant)
 return err
}

// CheckTenantSession rejects sessions issued before any disable transition.
// Refresh and reauthentication must never replace the original generation.
func CheckTenantSession(ctx context.Context, tenant string, issued int64) error {
 generation, err := TenantGeneration(ctx, tenant)
 if err != nil { return err }
 if generation != 0 && issued != generation { return contracts.ErrTenantDenied }
 return nil
}

// CheckAccountTenant derives scope from the stored account, not request input.
func CheckAccountTenant(ctx context.Context, tenant *string) error {
 if tenant == nil { return CheckTenant(ctx, "") }
 return CheckTenant(ctx, *tenant)
}

// CheckAccountGeneration validates the original generation of a credential
// against the tenant derived from its stored account.
func CheckAccountGeneration(ctx context.Context, tenant *string, issued int64) error {
 if tenant == nil { return CheckTenantSession(ctx, "", issued) }
 return CheckTenantSession(ctx, *tenant, issued)
}
