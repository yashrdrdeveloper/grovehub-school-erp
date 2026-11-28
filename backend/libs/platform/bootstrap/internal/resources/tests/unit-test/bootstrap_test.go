package bootstrap_test

import (
	"context"
	"io"
	"testing"

	"grovehub/libs/platform/bootstrap"
	"grovehub/libs/platform/command"
)

type strategy struct{}

func (strategy) Run([]string, io.Writer, io.Writer) int { return 0 }

func TestRegistration(t *testing.T) {
	called := false
	registration := bootstrap.Registration[string]{
		ServiceName: "example",
		RegisterCommands: func(db command.Database[string]) map[string]command.Command {
			if db != nil {
				t.Fatal("unexpected database")
			}
			return map[string]command.Command{"example": strategy{}}
		},
		MigrationPlan: bootstrap.Migration[string]{
			Target: bootstrap.AuditDatabase,
			Run: func(_ context.Context, value string) error {
				called = value == "database"
				return nil
			},
		},
	}
	var service bootstrap.Service[string] = registration
	if service.Name() != "example" || len(service.Commands(nil)) != 1 {
		t.Fatal("registration was not forwarded")
	}
	migration := service.Migration()
	if migration.Target != bootstrap.AuditDatabase || migration.Run(context.Background(), "database") != nil || !called {
		t.Fatal("migration was not forwarded")
	}
	if commands := (bootstrap.Registration[string]{}).Commands(nil); commands != nil {
		t.Fatal("a service without commands must be empty")
	}
}
