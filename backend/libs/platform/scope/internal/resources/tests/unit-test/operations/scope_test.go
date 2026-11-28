package operations_test

import (
	"context"
	"errors"
	"testing"

	"grovehub/libs/platform/scope"
)

func TestTargets(t *testing.T) {
	for _, tc := range []struct {
		target scope.Target
		ok bool
	}{
		{scope.Target{Level: scope.Platform}, true},
		{scope.Target{Level: scope.Platform, ProjectID: "p"}, false},
		{scope.Target{Level: scope.Project, ProjectID: "p"}, true},
		{scope.Target{Level: scope.Project, ProjectID: "p", TenantID: "t"}, false},
		{scope.Target{Level: scope.Project}, false},
		{scope.Target{Level: scope.Tenant, ProjectID: "p", TenantID: "t"}, true},
		{scope.Target{Level: scope.Tenant, ProjectID: "p"}, false},
		{scope.Target{Level: scope.Tenant, TenantID: "t"}, false},
		{scope.Target{Level: scope.Tenant, ProjectID: " p", TenantID: "t"}, false},
		{scope.Target{Level: scope.Tenant, ProjectID: "p", TenantID: "t\x00"}, false},
		{scope.Target{}, false},
	} {
		err := scope.ValidateTarget(tc.target)
		if (err == nil) != tc.ok || (!tc.ok && !errors.Is(err, scope.ErrInvalid)) {
			t.Fatalf("target=%+v error=%v", tc.target, err)
		}
	}
}

func TestExecution(t *testing.T) {
	target := scope.Target{Level: scope.Tenant, ProjectID: "target", TenantID: "tenant"}
	workload := scope.Execution{Capability: "audit.write", Target: target, CallerProjectID: "caller", PrincipalType: scope.Workload, PrincipalID: "w", WorkloadID: "w"}
	user := workload
	user.PrincipalType, user.PrincipalID, user.SessionID = scope.User, "u", "s"
	system := scope.Execution{Capability: "audit.write", Target: scope.Target{Level: scope.Platform}, PrincipalType: scope.System, PrincipalID: "system"}
	for _, value := range []scope.Execution{workload, user, system} {
		if err := scope.ValidateExecution(value); err != nil {
			t.Fatal(err)
		}
		ctx, err := scope.WithVerifiedExecution(context.Background(), value)
		if err != nil {
			t.Fatal(err)
		}
		if actual, ok := scope.From(ctx); !ok || actual != value {
			t.Fatal("context was not preserved")
		}
		if actual, err := scope.Require(ctx, value.Target, value.Capability); err != nil || actual != value {
			t.Fatal("authorized target denied")
		}
		if _, err := scope.Require(ctx, scope.Target{}, value.Capability); !errors.Is(err, scope.ErrDenied) {
			t.Fatal("invalid target allowed")
		}
	}
	for _, change := range []func(*scope.Execution){
		func(v *scope.Execution) { v.Target = scope.Target{} },
		func(v *scope.Execution) { v.PrincipalID = "" },
		func(v *scope.Execution) { v.Capability = "" },
		func(v *scope.Execution) { v.CallerProjectID = " bad" },
		func(v *scope.Execution) { v.PrincipalType = "unknown" },
		func(v *scope.Execution) { v.WorkloadID = "other" },
		func(v *scope.Execution) { v.SessionID = "s" },
		func(v *scope.Execution) { v.PrincipalType = scope.User },
		func(v *scope.Execution) { v.PrincipalType = scope.System },
	} {
		value := workload
		change(&value)
		if _, err := scope.WithVerifiedExecution(context.Background(), value); !errors.Is(err, scope.ErrInvalid) {
			t.Fatalf("accepted invalid execution: %+v", value)
		}
	}
	if _, err := scope.WithVerifiedExecution(nil, workload); !errors.Is(err, scope.ErrInvalid) {
		t.Fatal("nil context accepted")
	}
	if _, ok := scope.From(nil); ok {
		t.Fatal("nil context accepted")
	}
	if _, err := scope.Require(context.Background(), target, "audit.write"); !errors.Is(err, scope.ErrDenied) {
		t.Fatal("missing context accepted")
	}
	ctx, _ := scope.WithVerifiedExecution(context.Background(), workload)
	if _, err := scope.Require(ctx, target, "audit.read"); !errors.Is(err, scope.ErrDenied) {
		t.Fatal("different capability accepted")
	}
	for _, other := range []scope.Target{
		{Level: scope.Tenant, ProjectID: "other", TenantID: "tenant"},
		{Level: scope.Tenant, ProjectID: "target", TenantID: "other"},
		{Level: scope.Project, ProjectID: "target"},
		{Level: scope.Platform},
	} {
		if _, err := scope.Require(ctx, other, "audit.write"); !errors.Is(err, scope.ErrDenied) {
			t.Fatalf("cross-scope access accepted: %+v", other)
		}
	}
}

func TestGrants(t *testing.T) {
	target := scope.Target{Level: scope.Tenant, ProjectID: "target", TenantID: "tenant"}
	value := scope.Execution{Capability: "audit.write", Target: target, CallerProjectID: "caller", PrincipalType: scope.Workload, PrincipalID: "w", WorkloadID: "w"}
	grant := scope.Grant{Target: target, CallerProjectID: "caller", WorkloadID: "w", Capability: "audit.write"}
	if !scope.Allows(grant, value, "audit.write") {
		t.Fatal("valid workload denied")
	}
	for _, change := range []func(*scope.Grant){
		func(g *scope.Grant) { g.Target = scope.Target{} },
		func(g *scope.Grant) { g.Target.TenantID = "other" },
		func(g *scope.Grant) { g.Target.ProjectID = "other" },
		func(g *scope.Grant) { g.CallerProjectID = "other" },
		func(g *scope.Grant) { g.WorkloadID = "" },
		func(g *scope.Grant) { g.WorkloadID = "other" },
		func(g *scope.Grant) { g.Capability = "audit.read" },
	} {
		other := grant
		change(&other)
		if scope.Allows(other, value, "audit.write") {
			t.Fatalf("unexpected grant: %+v", other)
		}
	}
	if scope.Allows(grant, value, "") || scope.Allows(grant, scope.Execution{}, "audit.write") {
		t.Fatal("invalid capability or execution allowed")
	}
	value.PrincipalType, value.PrincipalID, value.SessionID = scope.User, "user", "session"
	if scope.Allows(grant, value, "audit.write") {
		t.Fatal("delegation not denied")
	}
	grant.AllowDelegation = true
	if !scope.Allows(grant, value, "audit.write") {
		t.Fatal("allowed delegation denied")
	}
	system := scope.Execution{Capability: "audit.write", Target: scope.Target{Level: scope.Platform}, PrincipalType: scope.System, PrincipalID: "system"}
	if scope.Allows(grant, system, "audit.write") {
		t.Fatal("system received workload grant")
	}
}
