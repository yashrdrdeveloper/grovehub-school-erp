package evaluate_test

import (
	"context"
	"testing"

	coregateway "grovehub/libs/core/gateway"
	"grovehub/services/identity/internal/app/constants/feature"
	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/entity"
	"grovehub/services/identity/internal/app/services/authorization/operations/evaluate"
)

func TestAllowed(t *testing.T) {
	principal := coregateway.Principal{AccountID: "account-1", TenantID: "school-1", SessionID: "session-1"}
	ctx := coregateway.WithPrincipal(context.Background(), principal)
	valid := func() evaluate.Snapshot {
		return evaluate.Snapshot{
			AccountID: "account-1",
			TenantID: "school-1",
			Status: entity.Active,
			Roles: []evaluate.RoleGrant{
				{RoleID: "primary", TenantID: "school-1", Active: true},
				{RoleID: "secondary", TenantID: "school-1", Active: true, Permissions: []permission.Code{permission.AccountSearch}},
			},
			Permissions: map[permission.Code]evaluate.PermissionState{permission.AccountSearch: {Active: true, Feature: feature.AccountManagement}},
			Features: map[feature.Code]evaluate.FeatureState{feature.AccountManagement: {Active: true, Enabled: true}},
		}
	}

	if !evaluate.Allowed(ctx, permission.AccountSearch, valid()) {
		t.Fatal("an active secondary role must grant its registered permission")
	}

	cases := map[string]func(*evaluate.Snapshot){
		"foreign tenant": func(s *evaluate.Snapshot) { s.TenantID = "school-2" },
		"foreign account": func(s *evaluate.Snapshot) { s.AccountID = "account-2" },
		"suspended account": func(s *evaluate.Snapshot) { s.Status = entity.Suspended },
		"unregistered permission": func(s *evaluate.Snapshot) { delete(s.Permissions, permission.AccountSearch) },
		"disabled permission": func(s *evaluate.Snapshot) { s.Permissions[permission.AccountSearch] = evaluate.PermissionState{Feature: feature.AccountManagement} },
		"missing feature association": func(s *evaluate.Snapshot) { s.Permissions[permission.AccountSearch] = evaluate.PermissionState{Active: true} },
		"disabled feature": func(s *evaluate.Snapshot) { s.Features[feature.AccountManagement] = evaluate.FeatureState{Active: true} },
		"inactive feature": func(s *evaluate.Snapshot) { s.Features[feature.AccountManagement] = evaluate.FeatureState{Enabled: true} },
		"missing entitlement": func(s *evaluate.Snapshot) { delete(s.Features, feature.AccountManagement) },
		"foreign role": func(s *evaluate.Snapshot) { s.Roles[1].TenantID = "school-2" },
		"inactive role": func(s *evaluate.Snapshot) { s.Roles[1].Active = false },
		"no matching grant": func(s *evaluate.Snapshot) { s.Roles[1].Permissions = nil },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			snapshot := valid()
			mutate(&snapshot)
			if evaluate.Allowed(ctx, permission.AccountSearch, snapshot) {
				t.Fatal("invalid snapshot must deny access")
			}
		})
	}
	if evaluate.Allowed(context.Background(), permission.AccountSearch, valid()) {
		t.Fatal("missing verified principal must deny access")
	}
	if evaluate.Allowed(ctx, permission.Code("UNKNOWN"), valid()) {
		t.Fatal("unknown capability must deny access")
	}
	self := valid()
	self.Roles[1].Permissions = []permission.Code{permission.AccountSelfView}
	self.Permissions[permission.AccountSelfView] = evaluate.PermissionState{Active: true}
	if !evaluate.Allowed(ctx, permission.AccountSelfView, self) {
		t.Fatal("registered featureless capability must not require an entitlement")
	}
	self.Permissions[permission.AccountSelfView] = evaluate.PermissionState{Active: true, Feature: feature.AccountManagement}
	if evaluate.Allowed(ctx, permission.AccountSelfView, self) {
		t.Fatal("unexpected feature association must deny access even when enabled")
	}
	principal.Admin = true
	if evaluate.Allowed(coregateway.WithPrincipal(context.Background(), principal), permission.AccountDeactivate, valid()) {
		t.Fatal("legacy admin flag must not bypass RBAC")
	}
}
