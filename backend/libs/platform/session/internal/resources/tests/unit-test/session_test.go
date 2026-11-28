package session_test

import (
	"context"
	"errors"
	"testing"

	"grovehub/libs/platform/session"
)

func TestPrincipalFacade(t *testing.T) {
	ctx := context.Background()

	for _, value := range []session.Principal{{}, {AccountID: "account"}, {SessionID: "session"}} {
		if _, ok := session.PrincipalFrom(session.WithPrincipal(ctx, value)); ok {
			t.Fatalf("accepted incomplete principal: %+v", value)
		}
	}

	if _, ok := session.PrincipalFrom(ctx); ok {
		t.Fatal("accepted absent principal")
	}

	want := session.Principal{
		AccountID: "account", SessionID: "session", TenantID: "school",
		PersonGID: "person", SessionGeneration: 3,
	}
	if got, ok := session.PrincipalFrom(session.WithPrincipal(ctx, want)); !ok || got != want {
		t.Fatalf("principal not preserved: %+v, %v", got, ok)
	}
}

func TestTenantFacade(t *testing.T) {
	ctx := context.Background()
	if session.ActiveTenantCheck(nil) != nil || session.WithTenantCheck(ctx, nil) != ctx {
		t.Fatal("absent policy should leave isolated context unchanged")
	}

	if generation, err := session.TenantGeneration(ctx, ""); generation != 0 || err != nil {
		t.Fatalf("isolated generation: %d, %v", generation, err)
	}
	if err := session.CheckTenantSession(ctx, "school", 1); err != nil {
		t.Fatal(err)
	}

	lookup := func(_ context.Context, id string) (session.TenantState, error) {
		switch id {
		case "active":
			return session.TenantState{Generation: 3, Active: true}, nil
		case "disabled":
			return session.TenantState{Generation: 3}, nil
		case "zero":
			return session.TenantState{Active: true}, nil
		default:
			return session.TenantState{}, errors.New("unavailable")
		}
	}
	scoped := session.WithTenantCheck(ctx, session.ActiveTenantCheck(lookup))

	for _, id := range []string{"", "disabled", "zero", "unknown"} {
		if err := session.CheckTenant(scoped, id); !errors.Is(err, session.ErrTenantDenied) {
			t.Fatalf("%q: expected tenant denial, got %v", id, err)
		}
	}

	active := "active"
	if err := session.CheckAccountTenant(scoped, &active); err != nil {
		t.Fatal(err)
	}
	if err := session.CheckAccountTenant(scoped, nil); !errors.Is(err, session.ErrTenantDenied) {
		t.Fatalf("missing account tenant: %v", err)
	}
	if generation, err := session.TenantGeneration(scoped, active); generation != 3 || err != nil {
		t.Fatalf("generation: %d, %v", generation, err)
	}

	for _, issued := range []int64{0, 1, 4} {
		if err := session.CheckAccountGeneration(scoped, &active, issued); !errors.Is(err, session.ErrTenantDenied) {
			t.Fatalf("accepted incorrect generation %d: %v", issued, err)
		}
	}
	if err := session.CheckAccountGeneration(scoped, &active, 3); err != nil {
		t.Fatal(err)
	}
	if err := session.CheckAccountGeneration(scoped, nil, 3); !errors.Is(err, session.ErrTenantDenied) {
		t.Fatalf("missing account generation tenant: %v", err)
	}
	if err := session.CheckTenantSession(scoped, "disabled", 3); !errors.Is(err, session.ErrTenantDenied) {
		t.Fatalf("disabled session: %v", err)
	}

	for _, generation := range []int64{0, -1} {
		invalid := session.WithTenantCheck(ctx, func(context.Context, string) (int64, error) {
			return generation, nil
		})
		if _, err := session.TenantGeneration(invalid, active); !errors.Is(err, session.ErrTenantDenied) {
			t.Fatalf("accepted nonpositive generation %d: %v", generation, err)
		}
	}
}

func TestSessionCookieContract(t *testing.T) {
	got := []string{session.AccessCookie, session.RefreshCookie, session.CSRFName, session.ResetCookie, session.SetupCookie}
	want := []string{"__Host-gh_access", "__Secure-gh_refresh", "gh_csrf", "__Secure-gh_reset", "__Secure-gh_setup"}
	for index := range want {
		if got[index] != want[index] {
			t.Fatalf("cookie %d = %q, want %q", index, got[index], want[index])
		}
	}
}
