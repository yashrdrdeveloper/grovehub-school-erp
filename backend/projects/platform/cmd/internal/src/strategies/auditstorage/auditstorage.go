package auditstorage

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"time"

	"grovehub/cmd/internal/src/command"
	"grovehub/cmd/internal/src/operations/database"
	"grovehub/services/audit"
)

// @Strategy
type strategy struct {
	database database.Database
}

func New(db database.Database) command.Command {
	return &strategy{database: db}
}

func (s *strategy) Run(args []string, output, diagnostics io.Writer) int {
	if err := s.run(args, output, diagnostics); err != nil {
		_, _ = fmt.Fprintln(diagnostics, err)
		return 1
	}

	return 0
}

func (s *strategy) run(args []string, output, diagnostics io.Writer) error {
	flags := flag.NewFlagSet("audit-storage", flag.ContinueOnError)
	flags.SetOutput(diagnostics)
	database := flags.String("database", "", "Exact expected Audit database name")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 || *database == "" {
		return errors.New("database is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db, closeDB, err := s.database.Audit(ctx, nil)
	if err != nil {
		return errors.New("could not connect to configured Audit database")
	}
	defer closeDB()

	tables, err := audit.InspectStorage(ctx, db, *database)
	if err != nil {
		return errors.New("Audit storage inspection failed; verify database name, schema and permissions")
	}

	return json.NewEncoder(output).Encode(tables)
}
