package rbacmigrate_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"grovehub/cmd/internal/src/strategies/rbacmigrate"
)

func TestRBACMigrate(t *testing.T) {
	for _, tc := range []struct {
		err  error
		code int
	}{
		{},
		{err: errors.New("failure"), code: 1},
	} {
		var output, diagnostics bytes.Buffer
		args := []string{"plan", "--tenant-id", "school"}
		strategy := rbacmigrate.NewWithRunner(func(got []string, out, diag io.Writer) error {
			if !reflect.DeepEqual(got, args) || out != &output || diag != &diagnostics {
				t.Fatal("operator arguments or streams were not forwarded")
			}

			return tc.err
		})

		if code := strategy.Run(args, &output, &diagnostics); code != tc.code {
			t.Fatalf("code=%d want=%d", code, tc.code)
		}
	}
}
