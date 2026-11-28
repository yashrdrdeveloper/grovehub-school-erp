package auditstorage_test

import (
	"context"
	"io"
	"reflect"
	"strings"
	"testing"

	"gorm.io/gorm"

	"grovehub/cmd/internal/resources/tests/helpers/fixture"
	"grovehub/cmd/internal/src/strategies/auditstorage"
	"grovehub/services/audit"
)

func TestAuditStorage(t *testing.T) {
	required := []string{"--database", "audit"}
	for _, tc := range []struct {
		name, stage string
		args        []string
		code        int
		closed      []string
	}{
		{name: "help", args: []string{"--help"}},
		{name: "bad flag", args: []string{"--bad"}, code: 1},
		{name: "missing database", code: 1},
		{name: "open failure", args: required, stage: "open", code: 1},
		{name: "inspection failure", args: required, stage: "inspect", code: 1, closed: []string{"audit"}},
		{name: "success", args: required, stage: "success", closed: []string{"audit"}},
		{name: "output failure", args: required, stage: "output", code: 1, closed: []string{"audit"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := fixture.New(t)
			if tc.stage == "open" {
				f.ExpectAudit(nil, fixture.Failure)
			} else if tc.stage != "" {
				f.ExpectAudit(nil, nil)
			}

			strategy := auditstorage.NewWithInspector(f.Mock, func(_ context.Context, db *gorm.DB, expected string) ([]audit.StorageTable, error) {
				if db != f.AuditDB || expected != "audit" {
					t.Fatal("incorrect audit scope")
				}
				if tc.stage == "inspect" {
					return nil, fixture.Failure
				}

				return []audit.StorageTable{{Name: "activities", TotalBytes: 123}}, nil
			})
			var output io.Writer = &f.Output
			if tc.stage == "output" {
				output = fixture.Writer{}
			}
			code := strategy.Run(tc.args, output, &f.Diagnostics)
			if code != tc.code || !reflect.DeepEqual(f.Closed, tc.closed) {
				t.Fatalf("code=%d cleanup=%v diagnostics=%s", code, f.Closed, f.Diagnostics.String())
			}
			if tc.stage == "success" && !strings.Contains(f.Output.String(), "activities") {
				t.Fatal("JSON result missing")
			}
		})
	}
}
