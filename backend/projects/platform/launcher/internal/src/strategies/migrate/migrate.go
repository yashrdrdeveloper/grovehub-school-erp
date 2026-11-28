package migrate

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"gorm.io/gorm"

	"grovehub/libs/platform/bootstrap"
	"grovehub/libs/platform/command"
)

type Target struct {
	Name      string
	Migration bootstrap.Migration[*gorm.DB]
}

// @Strategy
type strategy struct {
	database command.Database[*gorm.DB]
	targets  []Target
}

func New(db command.Database[*gorm.DB], targets []Target) command.Command {
	return &strategy{database: db, targets: targets}
}

func (s *strategy) Run(args []string, _ io.Writer, diagnostics io.Writer) int {
	names := make([]string, 0, len(s.targets)+1)
	selected := make([]Target, 0, len(s.targets))
	for _, target := range s.targets {
		names = append(names, target.Name)
		if len(args) == 1 && (args[0] == "all" || args[0] == target.Name) {
			selected = append(selected, target)
		}
	}

	if len(selected) == 0 {
		names = append(names, "all")
		_, _ = fmt.Fprintf(diagnostics, "usage: groovehub-platform migrate %s\n", strings.Join(names, "|"))
		return 2
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	var business *gorm.DB
	for _, target := range selected {
		if target.Migration.Target == bootstrap.AuditDatabase {
			if err := s.audit(ctx, business, target); err != nil {
				_, _ = fmt.Fprintln(diagnostics, err)
				return 1
			}
			continue
		}

		if business == nil {
			db, closeDB, err := s.database.Business(ctx)
			if err != nil {
				_, _ = fmt.Fprintf(diagnostics, "open migration database: %v\n", err)
				return 1
			}
			business = db
			defer closeDB()
		}

		if err := target.Migration.Run(ctx, business); err != nil {
			_, _ = fmt.Fprintf(diagnostics, "migrate %s: %v\n", target.Name, err)
			return 1
		}
	}

	return 0
}

func (s *strategy) audit(ctx context.Context, shared *gorm.DB, target Target) error {
	db, closeDB, err := s.database.Audit(ctx, shared)
	if err != nil {
		return fmt.Errorf("open Audit migration database: %w", err)
	}
	defer closeDB()

	if err := target.Migration.Run(ctx, db); err != nil {
		return fmt.Errorf("migrate %s: %w", target.Name, err)
	}

	return nil
}
