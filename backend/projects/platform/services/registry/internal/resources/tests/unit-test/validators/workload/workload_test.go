package workload_test

import (
	"encoding/base64"
	"strings"
	"testing"
	"time"

	"grovehub/libs/platform/scope"
	dto "grovehub/services/registry/internal/app/dto/workload"
	"grovehub/services/registry/internal/app/validators/workload"
)

func TestKeyValidity(t *testing.T) {
	now := time.Now().UTC()
	valid := dto.KeyInput{KeyID: "key-1", PublicKey: base64.RawURLEncoding.EncodeToString(make([]byte, 32)), ValidFrom: now, ValidUntil: now.Add(time.Hour)}
	if key, err := workload.ValidateKey(valid, now); err != nil || len(key) != 32 {
		t.Fatalf("valid public key encoding rejected: %v", err)
	}
	for _, change := range []func(*dto.KeyInput){
		func(v *dto.KeyInput) { v.KeyID = "bad key" },
		func(v *dto.KeyInput) { v.PublicKey += "=" },
		func(v *dto.KeyInput) { v.PublicKey = strings.Repeat("a", 42) },
		func(v *dto.KeyInput) { v.ValidFrom = time.Time{} },
		func(v *dto.KeyInput) { v.ValidUntil = now },
		func(v *dto.KeyInput) { v.ValidFrom = now.Add(2 * time.Hour) },
		func(v *dto.KeyInput) { v.ValidUntil = now.Add(workload.MaxKeyLifetime + time.Hour) },
	} {
		bad := valid
		change(&bad)
		if _, err := workload.ValidateKey(bad, now); err == nil {
			t.Fatalf("invalid key accepted: %+v", bad)
		}
	}
}

func TestGrantScope(t *testing.T) {
	const project = "10000000-0000-4000-8000-000000000001"
	const tenant = "20000000-0000-4000-8000-000000000001"
	for _, grant := range []dto.GrantInput{
		{Scope: scope.Platform, Capability: "registry.workload.create"},
		{Scope: scope.Project, TargetProjectID: project, Capability: "registry.project.read"},
		{Scope: scope.Tenant, TargetProjectID: project, TargetTenantID: tenant, Capability: "audit.read", AllowDelegation: true},
	} {
		if err := workload.ValidateGrant(grant); err != nil {
			t.Fatalf("valid scope rejected: %v", err)
		}
	}
	for _, grant := range []dto.GrantInput{
		{Scope: scope.Platform, Capability: "registry.workload.create", AllowDelegation: true},
		{Scope: scope.Platform, TargetProjectID: project, Capability: "audit.read"},
		{Scope: scope.Project, TargetProjectID: project, TargetTenantID: tenant, Capability: "audit.read"},
		{Scope: scope.Tenant, TargetProjectID: project, Capability: "audit.read"},
		{Scope: scope.Tenant, TargetProjectID: "bad", TargetTenantID: tenant, Capability: "audit.read"},
		{Scope: scope.Tenant, TargetProjectID: project, TargetTenantID: "bad", Capability: "audit.read"},
		{Scope: scope.Project, TargetProjectID: project, Capability: "bad capability"},
	} {
		if err := workload.ValidateGrant(grant); err == nil {
			t.Fatalf("invalid grant accepted: %+v", grant)
		}
	}
}
