package schoolgateway

import (
 "context"
 "errors"
 "testing"
)

func TestTenantSessionGeneration(t *testing.T) {
 generation := int64(1)
 active := true
 ctx := WithTenantCheck(context.Background(), func(_ context.Context, tenant string) (int64, error) {
  if tenant != "school" || !active { return 0, ErrTenantDenied }
  return generation, nil
 })
 issued, err := TenantGeneration(ctx, "school")
 if err != nil || issued != 1 { t.Fatal(issued, err) }
 if err := CheckTenantSession(ctx, "school", issued); err != nil { t.Fatal(err) }
 active = false
 generation++
 if err := CheckTenantSession(ctx, "school", issued); !errors.Is(err, ErrTenantDenied) { t.Fatal("disabled school accepted", err) }
 active = true
 for _, old := range []int64{0, issued, generation+1} {
  if err := CheckTenantSession(ctx, "school", old); !errors.Is(err, ErrTenantDenied) { t.Fatal("incorrect generation accepted", old, err) }
 }
 if err := CheckTenantSession(ctx, "school", generation); err != nil { t.Fatal("new session denied", err) }
 if err := CheckTenantSession(ctx, "other", generation); !errors.Is(err, ErrTenantDenied) { t.Fatal("wrong tenant accepted", err) }
}

func TestCheckAccountTenant(t *testing.T) {
 tenant := "stored-school"
 ctx := WithTenantCheck(context.Background(), func(_ context.Context, got string) (int64,error) {
  if got != tenant { t.Fatal("incorrect tenant",got) }
  return 0, ErrTenantDenied
 })
 if !errors.Is(CheckAccountTenant(ctx, &tenant), ErrTenantDenied) { t.Fatal("disabled account accepted") }
 if !errors.Is(CheckAccountTenant(ctx, nil), ErrTenantDenied) { t.Fatal("missing tenant accepted") }
}
