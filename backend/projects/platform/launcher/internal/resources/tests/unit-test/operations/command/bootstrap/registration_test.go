package bootstrap_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"gorm.io/gorm"

	"grovehub/launcher/internal/src/operations/command/bootstrap"
	"grovehub/launcher/internal/src/operations/command/registry"
	platformbootstrap "grovehub/libs/platform/bootstrap"
	"grovehub/libs/platform/command"
)

type strategy struct{}

func (strategy) Run(args []string, output, _ io.Writer) int {
	_, _ = io.WriteString(output, args[0])
	return 0
}

func registration(name, operation string, runner command.Command) platformbootstrap.Registration[*gorm.DB] {
	return platformbootstrap.Registration[*gorm.DB]{
		ServiceName: name,
		RegisterCommands: func(command.Database[*gorm.DB]) map[string]command.Command {
			return map[string]command.Command{operation: runner}
		},
	}
}

func TestServiceCommandRegistration(t *testing.T) {
	service := registration("example", "example-command", strategy{})
	app, err := bootstrap.NewWithServices(nil, nil, []platformbootstrap.Service[*gorm.DB]{service})
	if err != nil {
		t.Fatal(err)
	}

	var output bytes.Buffer
	if code := app.Run([]string{"example-command", "payload"}, &output, io.Discard); code != 0 || output.String() != "payload" {
		t.Fatalf("code=%d output=%q", code, output.String())
	}
}

func TestInvalidRegistration(t *testing.T) {
	for _, services := range [][]platformbootstrap.Service[*gorm.DB]{
		{registration("", "example", strategy{})},
		{registration("all", "example", strategy{})},
		{registration("example", "one", strategy{}), registration("example", "two", strategy{})},
		{registration("example", "serve", strategy{})},
		{registration("example", "migrate", strategy{})},
		{registration("example", "", strategy{})},
		{registration("example", "nil", nil)},
		{registration("one", "duplicate", strategy{}), registration("two", "duplicate", strategy{})},
	} {
		if app, err := bootstrap.NewWithServices(nil, nil, services); app != nil || err == nil {
			t.Fatalf("invalid registration accepted: %v", services)
		}
	}
}

func TestDispatchBootstrapFailure(t *testing.T) {
	var diagnostics bytes.Buffer
	code := bootstrap.Dispatch(nil, io.Discard, &diagnostics, func() (registry.Registry, error) {
		return nil, errors.New("invalid service registration")
	})
	if code != 1 || diagnostics.String() != "invalid service registration\n" {
		t.Fatalf("code=%d diagnostics=%q", code, diagnostics.String())
	}
}
