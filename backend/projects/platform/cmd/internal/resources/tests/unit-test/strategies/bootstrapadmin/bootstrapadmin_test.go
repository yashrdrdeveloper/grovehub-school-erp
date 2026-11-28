package bootstrapadmin_test

import (
	"bytes"
	"errors"
	"flag"
	"io"
	"reflect"
	"testing"

	"grovehub/cmd/internal/src/strategies/bootstrapadmin"
)

func TestBootstrapAdmin(t *testing.T) {
	for _, tc := range []struct {
		err  error
		code int
	}{
		{},
		{err: flag.ErrHelp},
		{err: errors.New("failure"), code: 1},
	} {
		var output, diagnostics bytes.Buffer
		args := []string{"--email", "operator@example.com"}
		strategy := bootstrapadmin.NewWithRunner(func(got []string, input io.Reader, out, diag io.Writer) error {
			if !reflect.DeepEqual(got, args) || input == nil || out != &output || diag != &diagnostics {
				t.Fatal("operator arguments or streams were not forwarded")
			}

			return tc.err
		})

		if code := strategy.Run(args, &output, &diagnostics); code != tc.code {
			t.Fatalf("code=%d want=%d", code, tc.code)
		}
	}
}
