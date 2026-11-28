package validator_test

import (
	"context"
	"testing"

	schoolgateway "grovehub/launcher/session"
	"grovehub/services/identity/internal/app/services/account/operations/roles"
)

func TestRoleCatalogRequiresAuthenticatedTenant(t *testing.T) {
	if _, err := roles.List(context.Background()); err == nil {
		t.Fatal("anonymous role catalog accepted")
	}
	noTenant := schoolgateway.WithPrincipal(context.Background(), schoolgateway.Principal{
		AccountID: "account",
		SessionID: "session",
		Admin:     true,
	})
	if _, err := roles.List(noTenant); err == nil {
		t.Fatal("tenantless role catalog accepted")
	}
	for _, admin := range []bool{false, true} {
		ctx := schoolgateway.WithPrincipal(context.Background(), schoolgateway.Principal{
			AccountID: "account",
			SessionID: "session",
			TenantID:  "tenant",
			Admin:     admin,
		})
		result, err := roles.List(ctx)
		if err != nil || len(result) == 0 {
			t.Fatalf("authenticated tenant role catalog failed: %v", err)
		}
	}
}
