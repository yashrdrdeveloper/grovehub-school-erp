package principal_test

import (
	"context"
	"testing"

	"grovehub/services/integration/client/principal"
)

func TestVerifiedPrincipal(t *testing.T) {
	ctx := context.Background()
	if _, ok := principal.From(ctx); ok {
		t.Fatal("missing principal accepted")
	}
	for _, value := range []principal.VerifiedPrincipal{{TenantID: "tenant"}, {ActorGID: "actor"}} {
		if _, ok := principal.From(principal.WithVerifiedPrincipal(ctx, value)); ok {
			t.Fatal("incomplete principal accepted")
		}
	}
	value := principal.VerifiedPrincipal{TenantID: "tenant", ActorGID: "actor"}
	if got, ok := principal.From(principal.WithVerifiedPrincipal(ctx, value)); !ok || got != value {
		t.Fatal("valid principal was not preserved")
	}
	if _, ok := principal.From(ctx); ok {
		t.Fatal("principal leaked to parent context")
	}
}
