package migrate

import (
	"context"
	"fmt"
	"io"
	"time"

	"gorm.io/gorm"

	"grovehub/cmd/internal/src/command"
	"grovehub/cmd/internal/src/operations/database"
	"grovehub/cmd/internal/src/operations/migration"
)

const usage = "usage: groovehub-platform migrate foundation|infrastructure|integration|identity|audit|all"

var order = [...]string{"foundation", "infrastructure", "integration", "identity", "audit"}

// @Strategy
type strategy struct {
	database database.Database
	migrations map[string]migration.Migration
}

func New(db database.Database, migrations map[string]migration.Migration) command.Command {
	return &strategy{database: db, migrations: migrations}
}

func (s *strategy) Run(args []string, _ io.Writer, diagnostics io.Writer) int {
	if len(args) != 1 {
		_, _ = fmt.Fprintln(diagnostics, usage)
		return 2
	}

	name := args[0]
	if name != "all" && !valid(name) {
		_, _ = fmt.Fprintln(diagnostics, usage)
		return 2
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	if name == "audit" {
		return s.audit(ctx, nil, diagnostics)
	}

	db, closeDB, err := s.database.Business(ctx)
	if err != nil {
		_, _ = fmt.Fprintf(diagnostics, "open migration database: %v\n", err)
		return 1
	}
	defer closeDB()

	for _, service := range order {
		if name != "all" && name != service {
			continue
		}

		if service == "audit" {
			return s.audit(ctx, db, diagnostics)
		}

		if err := s.migrations[service].Run(ctx, db); err != nil {
			_, _ = fmt.Fprintf(diagnostics, "migrate %s: %v\n", service, err)
			return 1
		}
	}

	return 0
}

func valid(name string) bool {
	for _, service := range order {
		if name == service {
			return true
		}
	}

	return false
}

func (s *strategy) audit(ctx context.Context, shared *gorm.DB, diagnostics io.Writer) int {
	db, closeDB, err := s.database.Audit(ctx, shared)
	if err != nil {
		_, _ = fmt.Fprintln(diagnostics, "open Audit migration database:", err)
		return 1
	}
	defer closeDB()

	if err := s.migrations["audit"].Run(ctx, db); err != nil {
		_, _ = fmt.Fprintln(diagnostics, "migrate audit:", err)
		return 1
	}

	return 0
}
