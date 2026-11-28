package config_test

import (
	"testing"

	config "grovehub/services/identity/internal/app/config/impl"
)

func TestApprovalExecutionConfig(t *testing.T) {
	t.Setenv("IDENTITY_RBAC_MODE", "strict")
	t.Setenv("IDENTITY_ALLOWED_ORIGIN", "")
	t.Setenv("IDENTITY_APPROVAL_EXECUTION_TENANT", "")
	cfg, err := config.Load()
	if err != nil || cfg.ApprovalExecutionTenant != "" {
		t.Fatalf("execution must default to disabled: %v", err)
	}
	t.Setenv("IDENTITY_APPROVAL_EXECUTION_TENANT", "school")
	cfg, err = config.Load()
	if err != nil || cfg.ApprovalExecutionTenant != "school" {
		t.Fatalf("exact tenant should be accepted: %v", err)
	}
	t.Setenv("IDENTITY_APPROVAL_EXECUTION_TENANT", "*")
	if _, err := config.Load(); err == nil {
		t.Fatal("wildcard must be rejected")
	}
	t.Setenv("IDENTITY_RBAC_MODE", "legacy-development")
	t.Setenv("IDENTITY_APPROVAL_EXECUTION_TENANT", "")
	t.Setenv("OTEL_DEPLOYMENT_ENVIRONMENT", "development")
	t.Setenv("IDENTITY_ALLOWED_ORIGIN", "https://school.example")
	if _, err := config.Load(); err == nil {
		t.Fatal("public legacy origin must be rejected")
	}
	t.Setenv("IDENTITY_ALLOWED_ORIGIN", "http://127.0.0.1:4173")
	if _, err := config.Load(); err != nil {
		t.Fatalf("explicit local development must be supported: %v", err)
	}
	t.Setenv("OTEL_DEPLOYMENT_ENVIRONMENT", "production")
	if _, err := config.Load(); err == nil {
		t.Fatal("production legacy mode must be rejected")
	}
}
