package gateway_test

import (
	. "grovehub/libs/core/gateway"

	"context"
	"crypto/ed25519"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestScopedWorkloadVerifier(t *testing.T) {
	public, private, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatal(err)
	}
	const route = "POST /foundation/restricted/persons/search"
	credentials := Ed25519Credentials{Keys: map[string]TrustedWorkloadKey{
		"school-1": {PublicKey: public, Issuer: "school", Subject: "school", Scopes: []string{"persons.search"}, AllowDelegation: true},
	}}
	verifier := ScopedWorkloadVerifier{Credentials: credentials, Audience: "platform", Routes: map[string]RoutePolicy{
		route: {Scope: "persons.search", RequireActor: true},
	}}
	valid := WorkloadClaims{Subject: "school", Audience: "platform", ExpiresAt: time.Now().Add(time.Minute), Scopes: []string{"persons.search"}, TenantID: "tenant", ActorGID: "actor"}
	sign := func(claims WorkloadClaims) string {
		t.Helper()
		token, err := SignWorkload("school-1", "school", private, claims)
		if err != nil {
			t.Fatal(err)
		}
		return token
	}
	token := sign(valid)
	wrongAudience := valid
	wrongAudience.Audience = "internal"
	wrongSubject := valid
	wrongSubject.Subject = "identity"
	wrongScope := valid
	wrongScope.Scopes = []string{"persons.write"}
	noActor := valid
	noActor.TenantID, noActor.ActorGID = "", ""
	parts := strings.Split(token, ".")
	parts[2] = "e30"
	cases := []struct {
		name string
		token string
		pattern string
		duplicate bool
		allow bool
	}{
		{"valid", token, route, false, true},
		{"missing", "", route, false, false},
		{"duplicate", token, route, true, false},
		{"unknown route", token, "POST /foundation/restricted/other", false, false},
		{"wrong audience", sign(wrongAudience), route, false, false},
		{"wrong subject", sign(wrongSubject), route, false, false},
		{"scope escalation", sign(wrongScope), route, false, false},
		{"missing actor", sign(noActor), route, false, false},
		{"tampered payload", strings.Join(parts, "."), route, false, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest("POST", "/foundation/restricted/persons/search", nil)
			request.Pattern = tc.pattern
			request.Header.Set("X-Actor-GID", "forged")
			if tc.token != "" {
				request.Header.Set("Authorization", "Bearer "+tc.token)
			}
			if tc.duplicate {
				request.Header.Add("Authorization", "Bearer "+tc.token)
			}
			ctx, err := verifier.VerifyWorkload(context.Background(), request)
			if !tc.allow {
				if err == nil || ctx != nil {
					t.Fatal("expected fail-closed denial")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			claims, ok := WorkloadFrom(ctx)
			if !ok || claims.ActorGID != "actor" || claims.Subject != "school" || claims.TenantID != "tenant" {
				t.Fatal("verified claims missing or incorrect")
			}
		})
	}
}

func TestSignWorkloadLifetime(t *testing.T) {
	_, private, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, lifetime := range []time.Duration{-time.Minute, 3*time.Minute} {
		_, err := SignWorkload("key", "school", private, WorkloadClaims{Subject: "school", Audience: "platform", ExpiresAt: time.Now().Add(lifetime)})
		if err == nil {
			t.Fatal("invalid lifetime accepted")
		}
	}
}
