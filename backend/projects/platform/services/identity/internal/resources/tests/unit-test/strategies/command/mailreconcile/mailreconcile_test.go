package mailreconcile_test

import (
	"context"
	"io"
	"reflect"
	"strings"
	"testing"

	"gorm.io/gorm"

	"grovehub/services/identity/internal/resources/tests/helpers/operatorfixture"
	"grovehub/services/identity/internal/app/strategies/command/mailreconcile"
	"grovehub/services/integration/client/mail"
)

func TestMailReconcile(t *testing.T) {
	required := []string{"--tenant-id", "school", "--database", "business"}
	for _, tc := range []struct {
		name, stage, dsn string
		args             []string
		code             int
		closed           []string
	}{
		{name: "help", args: []string{"--help"}},
		{name: "bad flag", args: []string{"--bad"}, code: 1},
		{name: "missing scope", code: 1},
		{name: "invalid tenant", args: []string{"--tenant-id", " school ", "--database", "business"}, code: 1},
		{name: "missing DSN", args: required, code: 1},
		{name: "open failure", args: required, dsn: "configured", stage: "open", code: 1},
		{name: "name lookup failure", args: required, dsn: "configured", stage: "name-error", code: 1, closed: []string{"business"}},
		{name: "name mismatch", args: required, dsn: "configured", stage: "name", code: 1, closed: []string{"business"}},
		{name: "search failure", args: required, dsn: "configured", stage: "search", code: 1, closed: []string{"business"}},
		{name: "default wiring", args: append(append([]string(nil), required...), "--limit", "0"), dsn: "configured", stage: "default", code: 1, closed: []string{"business"}},
		{name: "success", args: required, dsn: "configured", stage: "success", closed: []string{"business"}},
		{name: "output failure", args: required, dsn: "configured", stage: "output", code: 1, closed: []string{"business"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("DATABASE_URL", tc.dsn)
			f := fixture.New(t)
			if tc.stage == "open" {
				f.ExpectBusiness(fixture.Failure)
			} else if tc.stage != "" {
				f.ExpectBusiness(nil)
				switch tc.stage {
				case "name-error":
					f.ExpectName(f.BusinessDB, "", fixture.Failure)
				case "name":
					f.ExpectName(f.BusinessDB, "unexpected", nil)
				default:
					f.ExpectName(f.BusinessDB, "business", nil)
				}
			}

			strategy := mailreconcile.NewWithSearch(f.Mock, func(_ context.Context, db *gorm.DB, input mail.SearchRequest) (mail.Report, error) {
				if db != f.BusinessDB || input.TenantID != "school" || input.Limit != 100 || input.After != "" {
					t.Fatal("incorrect query scope")
				}
				if tc.stage == "search" {
					return mail.Report{}, fixture.Failure
				}

				return mail.Report{TenantID: input.TenantID, Rows: []mail.Row{}}, nil
			})
			if tc.stage == "default" {
				strategy = mailreconcile.New(f.Mock)
			}
			var output io.Writer = &f.Output
			if tc.stage == "output" {
				output = fixture.Writer{}
			}
			code := strategy.Run(tc.args, output, &f.Diagnostics)
			if code != tc.code || !reflect.DeepEqual(f.Closed, tc.closed) {
				t.Fatalf("code=%d cleanup=%v diagnostics=%s", code, f.Closed, f.Diagnostics.String())
			}
			if tc.stage == "success" && !strings.Contains(f.Output.String(), "school") {
				t.Fatal("JSON result missing")
			}
		})
	}
}
