package auditreplay_test

import (
	"context"
	"io"
	"reflect"
	"strings"
	"testing"

	"gorm.io/gorm"

	"grovehub/cmd/internal/resources/tests/helpers/fixture"
	"grovehub/cmd/internal/src/strategies/auditreplay"
	"grovehub/services/audit"
)

func TestAuditReplay(t *testing.T) {
	required := []string{"--tenant-id", "school", "--event-id", "event", "--database", "business"}
	for _, tc := range []struct {
		name, stage, auditURL, diagnostic string
		args                              []string
		code                              int
		closed                            []string
	}{
		{name: "help", args: []string{"--help"}},
		{name: "bad flag", args: []string{"--bad"}, code: 1},
		{name: "missing required", code: 1, diagnostic: "required"},
		{name: "apply requires operator", args: appendCopy(required, "--apply"), code: 1, diagnostic: "operator"},
		{name: "separate audit name required", args: required, auditURL: "configured", code: 1, diagnostic: "audit-database is required"},
		{name: "business open failure", args: required, stage: "business-open", code: 1, diagnostic: "connect"},
		{name: "business name failure", args: required, stage: "business-name-error", code: 1, closed: []string{"business"}},
		{name: "business name mismatch", args: required, stage: "business-name", code: 1, closed: []string{"business"}},
		{name: "audit open failure", args: required, stage: "audit-open", code: 1, closed: []string{"business"}},
		{name: "audit name failure", args: required, stage: "audit-name-error", code: 1, closed: []string{"audit", "business"}},
		{name: "audit name mismatch", args: required, stage: "audit-name", code: 1, closed: []string{"audit", "business"}},
		{name: "replay failure", args: required, stage: "replay-error", code: 1, closed: []string{"audit", "business"}},
		{name: "default wiring", args: required, stage: "default", code: 1, closed: []string{"audit", "business"}},
		{name: "success preview", args: required, stage: "success", closed: []string{"audit", "business"}},
		{name: "success apply", args: appendCopy(required, "--audit-database", "audit", "--operator", "operator", "--apply"), auditURL: "configured", stage: "apply", closed: []string{"audit", "business"}},
		{name: "output failure", args: required, stage: "output", code: 1, closed: []string{"audit", "business"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("AUDIT_DATABASE_URL", tc.auditURL)
			f := fixture.New(t)
			if tc.stage != "" {
				configure(t, f, tc.stage)
			}

			strategy := auditreplay.NewWithReplay(f.Mock, func(_ context.Context, auditDB, businessDB *gorm.DB, request audit.ReplayRequest) (audit.ReplayResult, error) {
				if auditDB != f.AuditDB || businessDB != f.BusinessDB || request.TenantID != "school" || request.EventID != "event" {
					t.Fatal("request or database scope not forwarded")
				}
				if request.Apply != (tc.stage == "apply") || (request.Apply && request.OperatorLabel != "operator") {
					t.Fatal("operator authorization not forwarded")
				}
				if tc.stage == "replay-error" {
					return audit.ReplayResult{}, fixture.Failure
				}

				return audit.ReplayResult{EventID: request.EventID, Applied: request.Apply}, nil
			})
			if tc.stage == "default" {
				strategy = auditreplay.New(f.Mock)
			}
			var output io.Writer = &f.Output
			if tc.stage == "output" {
				output = fixture.Writer{}
			}

			code := strategy.Run(tc.args, output, &f.Diagnostics)
			if code != tc.code || !reflect.DeepEqual(f.Closed, tc.closed) {
				t.Fatalf("code=%d cleanup=%v diagnostics=%s", code, f.Closed, f.Diagnostics.String())
			}
			if tc.diagnostic != "" && !strings.Contains(f.Diagnostics.String(), tc.diagnostic) {
				t.Fatal(f.Diagnostics.String())
			}
			if (tc.stage == "success" || tc.stage == "apply") && !strings.Contains(f.Output.String(), "event") {
				t.Fatal("JSON result missing")
			}
		})
	}
}

func appendCopy(base []string, extra ...string) []string {
	return append(append([]string(nil), base...), extra...)
}

func configure(t *testing.T, f *fixture.Operator, stage string) {
	t.Helper()
	if stage == "business-open" {
		f.ExpectBusiness(fixture.Failure)
		return
	}
	f.ExpectBusiness(nil)
	switch stage {
	case "business-name-error":
		f.ExpectName(f.BusinessDB, "", fixture.Failure)
		return
	case "business-name":
		f.ExpectName(f.BusinessDB, "unexpected", nil)
		return
	}
	f.ExpectName(f.BusinessDB, "business", nil)
	if stage == "audit-open" {
		f.ExpectAudit(f.BusinessDB, fixture.Failure)
		return
	}
	f.ExpectAudit(f.BusinessDB, nil)
	switch stage {
	case "audit-name-error":
		f.ExpectName(f.AuditDB, "", fixture.Failure)
	case "audit-name":
		f.ExpectName(f.AuditDB, "unexpected", nil)
	case "apply":
		f.ExpectName(f.AuditDB, "audit", nil)
	default:
		f.ExpectName(f.AuditDB, "business", nil)
	}
}
