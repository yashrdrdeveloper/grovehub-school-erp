package principal_test

import (
	"context"
	"testing"

	"grovehub/libs/platform/session/internal/src/contracts"
	"grovehub/libs/platform/session/internal/src/operations/principal"
)

func TestPrincipalContext(t *testing.T) {
	ctx := context.Background()
	if _, ok := principal.PrincipalFrom(ctx); ok {
		t.Fatal("absent principal accepted")
	}
	for _, value := range []contracts.Principal{{AccountID: "a"}, {SessionID: "s"}, {AccountID: "a", SessionID: "s"}} {
		got, ok := principal.PrincipalFrom(principal.WithPrincipal(ctx, value))
		want := value.AccountID != "" && value.SessionID != ""
		if got != value || ok != want {
			t.Fatalf("principal = %+v, %v; expected %+v, %v", got, ok, value, want)
		}
	}
}
