package validator_test

import (
	"context"
	"errors"
	"testing"

	coregateway "grovehub/libs/core/gateway"
	"grovehub/services/identity/internal/app/services/security/principal"
)

func TestRequireTenant(t *testing.T) {
	for _, tc := range []struct {
		name string
		actor coregateway.Principal
		want error
	}{
		{name: "anonymous", want: principal.ErrAuthenticationRequired},
		{name: "missing session", actor: coregateway.Principal{AccountID: "account", TenantID: "tenant"}, want: principal.ErrAuthenticationRequired},
		{name: "missing tenant", actor: coregateway.Principal{AccountID: "account", SessionID: "session", Admin: true}, want: principal.ErrTenantForbidden},
		{name: "member", actor: coregateway.Principal{AccountID: "account", SessionID: "session", TenantID: "tenant"}},
		{name: "admin", actor: coregateway.Principal{AccountID: "account", SessionID: "session", TenantID: "tenant", Admin: true}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := coregateway.WithPrincipal(context.Background(), tc.actor)
			actor, err := principal.RequireTenant(ctx)
			if !errors.Is(err, tc.want) {
				t.Fatalf("unexpected error: %v; want %v", err, tc.want)
			}
			if tc.want == nil && actor.TenantID != tc.actor.TenantID {
				t.Fatal("tenant identity was not preserved")
			}
		})
	}
}
