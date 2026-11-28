package bootstrap_test

import (
	"bytes"
	"testing"

	"grovehub/cmd/internal/src/bootstrap"
)

func TestRegisteredCommandsRejectInvalidInput(t *testing.T) {
	app := bootstrap.New()
	for _, tc := range []struct {
		args []string
		code int
	}{
		{[]string{"unknown"}, 2},
		{[]string{"serve", "extra"}, 2},
		{[]string{"migrate"}, 2},
		{[]string{"audit-replay"}, 1},
		{[]string{"audit-storage"}, 1},
		{[]string{"identity-bootstrap-admin"}, 1},
		{[]string{"identity-mail-reconcile"}, 1},
		{[]string{"identity-rbac-migrate"}, 1},
	} {
		var output, diagnostics bytes.Buffer
		if code := app.Run(tc.args, &output, &diagnostics); code != tc.code || diagnostics.Len() == 0 {
			t.Fatalf("%v: code=%d diagnostics=%q", tc.args, code, diagnostics.String())
		}
	}
}
