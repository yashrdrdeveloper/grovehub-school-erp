package contracts_test

import (
	"context"
	"testing"

	"grovehub/libs/platform/session/internal/src/contracts"
)

type verifier struct{ principal contracts.Principal }

func (v verifier) VerifyAccess(context.Context, string) (contracts.Principal, error) {
	return v.principal, nil
}

func (v verifier) AuthorizeProtectedRoute(context.Context, string) error {
	return nil
}

var (
	_ contracts.AccessVerifier           = verifier{}
	_ contracts.ProtectedRouteAuthorizer = verifier{}
)

func TestPrincipalContracts(t *testing.T) {
	want := contracts.Principal{AccountID: "account", SessionID: "session", TenantID: "school"}
	var access contracts.AccessVerifier = verifier{principal: want}
	got, err := access.VerifyAccess(context.Background(), "secret")
	if err != nil || got != want {
		t.Fatalf("principal contract: %+v, %v", got, err)
	}
}
