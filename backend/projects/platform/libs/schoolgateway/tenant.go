package schoolgateway

import (
 "context"
 "errors"
)

// TenantCheck returns the current session generation of an active school.
type TenantCheck func(context.Context, string) (int64, error)
type tenantCheckKey struct{}
var ErrTenantDenied = errors.New("tenant access denied")

func WithTenantCheck(ctx context.Context, check TenantCheck) context.Context {
 if check == nil { return ctx }
 return context.WithValue(ctx, tenantCheckKey{}, check)
}

// TenantGeneration is called only with trusted tenant identity. Zero is
// reserved for isolated module/test use where no registry policy is configured.
func TenantGeneration(ctx context.Context, tenant string) (int64, error) {
 check, _ := ctx.Value(tenantCheckKey{}).(TenantCheck)
 if check == nil { return 0, nil }
 if tenant == "" { return 0, ErrTenantDenied }
 generation, err := check(ctx, tenant)
 if err != nil || generation < 1 { return 0, ErrTenantDenied }
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
 if generation != 0 && issued != generation { return ErrTenantDenied }
 return nil
}
