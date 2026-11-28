package launcher_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"grovehub/launcher"
)

func TestExecuteUnknownCommand(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("AUDIT_DATABASE_URL", "")
	var diagnostics bytes.Buffer
	code := launcher.Execute([]string{"unknown"}, io.Discard, &diagnostics)
	if code != 2 || !strings.Contains(diagnostics.String(), "identity-bootstrap-admin") {
		t.Fatalf("code=%d diagnostics=%q", code, diagnostics.String())
	}
}
