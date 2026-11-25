package validator_test

import (
	"context"
	"testing"

	coregateway "grovehub/libs/core/gateway"
	"grovehub/services/identity/internal/app/services/account/operations/roles"
)

func TestRoleCatalogRequiresAuthenticatedTenant(t *testing.T) {
	if _, err := roles.List(context.Background()); err == nil {
		t.Fatal("anonymous role catalog accepted")
	}
	noTenant := coregateway.WithPrincipal(context.Background(), coregateway.Principal{
		AccountID: "account",
		SessionID: "session",
		Admin: true,
	})
	if _, err := roles.List(noTenant); err == nil {
		t.Fatal("tenantless role catalog accepted")
	}
	for _, admin := range []bool{false, true} {
		ctx := coregateway.WithPrincipal(context.Background(), coregateway.Principal{
			AccountID: "account",
			SessionID: "session",
			TenantID: "tenant",
			Admin: admin,
		})
		result, err := roles.List(ctx)
		if err != nil || len(result) == 0 {
			t.Fatalf("authenticated tenant role catalog failed: %v", err)
		}
	}
}
