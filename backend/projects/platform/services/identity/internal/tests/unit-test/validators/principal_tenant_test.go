package validator_test

import (
	"context"
	"errors"
	"testing"

	schoolgateway "grovehub/launcher/session"
	"grovehub/services/identity/internal/app/services/security/principal"
)

func TestRequireTenant(t *testing.T) {
	for _, tc := range []struct {
		name  string
		actor schoolgateway.Principal
		want  error
	}{
		{name: "anonymous", want: principal.ErrAuthenticationRequired},
		{name: "missing session", actor: schoolgateway.Principal{AccountID: "account", TenantID: "tenant"}, want: principal.ErrAuthenticationRequired},
		{name: "missing tenant", actor: schoolgateway.Principal{AccountID: "account", SessionID: "session", Admin: true}, want: principal.ErrTenantForbidden},
		{name: "member", actor: schoolgateway.Principal{AccountID: "account", SessionID: "session", TenantID: "tenant"}},
		{name: "admin", actor: schoolgateway.Principal{AccountID: "account", SessionID: "session", TenantID: "tenant", Admin: true}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := schoolgateway.WithPrincipal(context.Background(), tc.actor)
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
