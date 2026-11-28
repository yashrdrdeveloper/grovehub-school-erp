package bootstrap

import (
	"fmt"
	"io"
	"log/slog"

	"gorm.io/gorm"

	databaseimpl "grovehub/launcher/internal/src/operations/command/database/impl"
	"grovehub/launcher/internal/src/operations/command/registry"
	registryimpl "grovehub/launcher/internal/src/operations/command/registry/impl"
	"grovehub/launcher/internal/src/strategies/migrate"
	"grovehub/launcher/internal/src/strategies/serve"
	platformbootstrap "grovehub/libs/platform/bootstrap"
	"grovehub/libs/platform/command"
	"grovehub/services/audit"
	"grovehub/services/foundation"
	"grovehub/services/identity"
	"grovehub/services/infrastructure"
	"grovehub/services/integration"
)

// New composes service bootstraps without initializing their HTTP runtimes.
func New(run func(*slog.Logger) error) (registry.Registry, error) {
	return NewWithServices(databaseimpl.New(), serve.New(run), []platformbootstrap.Service[*gorm.DB]{
		foundation.Bootstrap(),
		infrastructure.Bootstrap(),
		integration.Bootstrap(),
		identity.Bootstrap(),
		audit.Bootstrap(),
	})
}

// NewWithServices preserves bootstrap order for migrations and rejects collisions.
func NewWithServices(db command.Database[*gorm.DB], server command.Command, services []platformbootstrap.Service[*gorm.DB]) (registry.Registry, error) {
	commands := map[string]command.Command{
		"serve":   server,
		"migrate": nil,
	}
	targets := make([]migrate.Target, 0, len(services))
	names := make(map[string]bool, len(services))

	for _, service := range services {
		name := service.Name()
		if name == "" || name == "all" || names[name] {
			return nil, fmt.Errorf("invalid or duplicate service bootstrap: %q", name)
		}
		names[name] = true
		migration := service.Migration()
		if migration.Run != nil {
			targets = append(targets, migrate.Target{Name: name, Migration: migration})
		}

		for name, strategy := range service.Commands(db) {
			if _, exists := commands[name]; name == "" || strategy == nil || exists {
				return nil, fmt.Errorf("invalid or duplicate command: %q", name)
			}
			commands[name] = strategy
		}
	}

	commands["migrate"] = migrate.New(db, targets)
	return registryimpl.New(commands), nil
}

// Dispatch converts bootstrap failures into the CLI failure code.
func Dispatch(args []string, output, diagnostics io.Writer, build func() (registry.Registry, error)) int {
	app, err := build()
	if err != nil {
		_, _ = fmt.Fprintln(diagnostics, err)
		return 1
	}

	return app.Run(args, output, diagnostics)
}
