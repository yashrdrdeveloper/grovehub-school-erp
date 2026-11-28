package gateway_test

import (
	"context"
	"errors"
	ginkgo "github.com/onsi/ginkgo/v2"

	schoolgateway "grovehub/internal/session"
	internalgateway "grovehub/launcher/internal/src/services/gateway/operations/inprocess"
	internalimpl "grovehub/launcher/internal/src/services/gateway/operations/inprocess/impl"
)

type allowOperation struct{}

func (allowOperation) Authorize(_ context.Context, principal schoolgateway.Principal, operation string) error {
	if principal.AccountID != "account" || operation != "account.get" {
		return errors.New("operation denied")
	}
	return nil
}

func verifyInternalGatewayRejectsUnknownPrincipal() {
	t := ginkgo.GinkgoT()
	gateway := internalimpl.New(allowOperation{})
	called := false

	_, err := internalgateway.Invoke(context.Background(), gateway, "account.get", func(context.Context) (string, error) {
		called = true
		return "private", nil
	})
	if err == nil || called {
		t.Fatal("unknown caller was allowed to invoke an internal operation")
	}
}

func verifyInternalGatewayKeepsTypedResult() {
	t := ginkgo.GinkgoT()
	gateway := internalimpl.New(allowOperation{})
	ctx := schoolgateway.WithPrincipal(context.Background(), schoolgateway.Principal{
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

func verifyInternalGatewayFailsClosedWithoutAuthorizer() {
	t := ginkgo.GinkgoT()
	gateway := internalimpl.New(nil)
	ctx := schoolgateway.WithPrincipal(context.Background(), schoolgateway.Principal{
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

// @Testing suite=LauncherInProcess
func InternalGatewaySpecs() {
	ginkgo.Describe("InternalGateway", func() {
		ginkgo.It("Internal Gateway Rejects Unknown Principal", verifyInternalGatewayRejectsUnknownPrincipal)
		ginkgo.It("Internal Gateway Keeps Typed Result", verifyInternalGatewayKeepsTypedResult)
		ginkgo.It("Internal Gateway Fails Closed Without Authorizer", verifyInternalGatewayFailsClosedWithoutAuthorizer)
	})
}
