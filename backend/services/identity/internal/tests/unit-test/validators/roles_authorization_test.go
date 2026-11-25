package validator_test

import (
	"context"
	"testing"

	coregateway "grovehub/libs/core/gateway"
	"grovehub/services/identity/internal/app/services/account/operations/roles"
)

func TestRoleCatalogRequiresTenantAdministrator(t *testing.T) {
	if _, err := roles.List(context.Background()); err == nil {
		t.Fatal("anonymous role catalog accepted")
	}
	ctx := coregateway.WithPrincipal(context.Background(), coregateway.Principal{
		AccountID: "account",
		SessionID: "session",
		TenantID: "tenant",
	})
	if _, err := roles.List(ctx); err == nil {
		t.Fatal("non-admin role catalog accepted")
	}
	ctx = coregateway.WithPrincipal(context.Background(), coregateway.Principal{
		AccountID: "account",
		SessionID: "session",
		TenantID: "tenant",
		Admin: true,
	})
	result, err := roles.List(ctx)
	if err != nil || len(result) == 0 {
		t.Fatalf("authorized role catalog failed: %v", err)
	}
}
