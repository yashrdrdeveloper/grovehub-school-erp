package serve_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	coretelemetry "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/telemetry"
	"grovehub/cmd/internal/src/strategies/serve"
)

func TestServe(t *testing.T) {
	for _, tc := range []struct {
		name                              string
		args                              []string
		loggerError, runError, closeError bool
		code                              int
		started                           bool
	}{
		{name: "unexpected args", args: []string{"extra"}, code: 2},
		{name: "logger failure", loggerError: true, code: 1},
		{name: "success", started: true},
		{name: "run failure", runError: true, code: 1, started: true},
		{name: "telemetry cleanup failure", closeError: true, started: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var output, diagnostics bytes.Buffer
			started, closed := false, false
			dependencies := serve.Dependencies{
				Logger: func(w io.Writer) (*slog.Logger, error) {
					if tc.loggerError {
						return nil, errors.New("logger")
					}

					return slog.New(slog.NewTextHandler(w, nil)), nil
				},
				Initialize: func(context.Context, *slog.Logger) coretelemetry.Shutdown {
					return func(context.Context) error {
						closed = true
						if tc.closeError {
							return errors.New("shutdown")
						}

						return nil
					}
				},
				Run: func(*slog.Logger) error {
					started = true
					if tc.runError {
						return errors.New("server")
					}

					return nil
				},
			}

			code := serve.NewWithDependencies(dependencies).Run(tc.args, &output, &diagnostics)
			if code != tc.code || started != tc.started || closed != tc.started {
				t.Fatalf("code=%d started=%v closed=%v", code, started, closed)
			}
			if tc.closeError && !bytes.Contains(output.Bytes(), []byte("OpenTelemetry shutdown failed")) {
				t.Fatal("missing cleanup diagnostic")
			}
		})
	}
}
