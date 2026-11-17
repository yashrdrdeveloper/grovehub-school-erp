package gateway_test

import (
	"context"
	"errors"
	"testing"

	coregateway "grovehub/libs/core/gateway"
	internalgateway "grovehub/launcher/services/gateway/operations/internal"
)

type allowOperation struct{}

func (allowOperation) Authorize(_ context.Context, principal coregateway.Principal, operation string) error {
	if principal.AccountID != "account" || operation != "account.get" {
		return errors.New("operation denied")
	}
	return nil
}

func TestInternalGatewayRejectsUnknownPrincipal(t *testing.T) {
	gateway := internalgateway.New(allowOperation{})
	called := false

	_, err := internalgateway.Invoke(context.Background(), gateway, "account.get", func(context.Context) (string, error) {
		called = true
		return "private", nil
	})
	if err == nil || called {
		t.Fatal("unknown caller was allowed to invoke an internal operation")
	}
}

func TestInternalGatewayKeepsTypedResult(t *testing.T) {
	gateway := internalgateway.New(allowOperation{})
	ctx := coregateway.WithPrincipal(context.Background(), coregateway.Principal{
		AccountID: "account",
		SessionID: "session",
	})

	result, err := internalgateway.Invoke(ctx, gateway, "account.get", func(context.Context) (string, error) {
		return "account-data", nil
	})
	if err != nil || result != "account-data" {
		t.Fatalf("typed invocation returned %q: %v", result, err)
	}
}

func TestInternalGatewayFailsClosedWithoutAuthorizer(t *testing.T) {
	gateway := internalgateway.New(nil)
	ctx := coregateway.WithPrincipal(context.Background(), coregateway.Principal{
		AccountID: "account",
		SessionID: "session",
	})
	_, err := internalgateway.Invoke(ctx, gateway, "account.get", func(context.Context) (int, error) {
		t.Fatal("unauthorized operation was invoked")
		return 0, nil
	})
	if err == nil {
		t.Fatal("missing authorizer was accepted")
	}
}
