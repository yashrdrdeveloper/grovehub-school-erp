package tenant_test

import (
	"context"
	"errors"
	"testing"

	"grovehub/libs/platform/session/internal/src/contracts"
	"grovehub/libs/platform/session/internal/src/operations/tenant"
)

func TestTenantPolicy(t *testing.T) {
	ctx := context.Background()
	if tenant.ActiveTenantCheck(nil) != nil || tenant.WithTenantCheck(ctx, nil) != ctx {
		t.Fatal("absent policy changed isolated context")
	}
	if n, err := tenant.TenantGeneration(ctx, ""); n != 0 || err != nil {
		t.Fatalf("isolated generation: %d, %v", n, err)
	}
	if err := tenant.CheckTenantSession(ctx, "school", 1); err != nil {
		t.Fatal(err)
	}

	generation := int64(1)
	active := true
	check := tenant.ActiveTenantCheck(func(_ context.Context, id string) (contracts.TenantState, error) {
		if id != "school" {
			return contracts.TenantState{}, errors.New("unknown tenant")
		}
		return contracts.TenantState{Generation: generation, Active: active}, nil
	})
	ctx = tenant.WithTenantCheck(ctx, check)
	id := "school"
	if err := tenant.CheckAccountTenant(ctx, &id); err != nil {
		t.Fatal(err)
	}
	if err := tenant.CheckAccountGeneration(ctx, &id, 1); err != nil {
		t.Fatal(err)
	}

	active = false
	generation++
	if err := tenant.CheckTenantSession(ctx, id, 1); !errors.Is(err, contracts.ErrTenantDenied) {
		t.Fatalf("disabled tenant accepted: %v", err)
	}
	active = true
	if err := tenant.CheckTenantSession(ctx, id, 1); !errors.Is(err, contracts.ErrTenantDenied) {
		t.Fatalf("stale credential restored: %v", err)
	}
	if err := tenant.CheckTenantSession(ctx, id, 2); err != nil {
		t.Fatal(err)
	}
	if err := tenant.CheckTenant(ctx, "unknown"); !errors.Is(err, contracts.ErrTenantDenied) {
		t.Fatalf("unknown tenant accepted: %v", err)
	}
	if err := tenant.CheckAccountTenant(ctx, nil); !errors.Is(err, contracts.ErrTenantDenied) {
		t.Fatalf("missing tenant accepted: %v", err)
	}
	if err := tenant.CheckAccountGeneration(ctx, nil, 2); !errors.Is(err, contracts.ErrTenantDenied) {
		t.Fatalf("missing tenant generation accepted: %v", err)
	}

	generation = 0
	if err := tenant.CheckTenant(ctx, id); !errors.Is(err, contracts.ErrTenantDenied) {
		t.Fatalf("zero lookup generation accepted: %v", err)
	}
	ctx = tenant.WithTenantCheck(ctx, func(context.Context, string) (int64, error) {
		return 0, nil
	})
	if err := tenant.CheckTenant(ctx, id); !errors.Is(err, contracts.ErrTenantDenied) {
		t.Fatalf("zero policy generation accepted: %v", err)
	}
}
