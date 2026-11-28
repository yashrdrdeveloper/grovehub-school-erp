package mailreconcile

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"grovehub/cmd/internal/src/command"
	"grovehub/cmd/internal/src/operations/database"
	"grovehub/services/identity"
	"grovehub/services/integration"
	"grovehub/services/integration/client/mail"
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
		_, _ = fmt.Fprintln(diagnostics, "identity-mail-reconcile:", err)
		return 1
	}

	return 0
}

func (s *strategy) run(args []string, output, diagnostics io.Writer) error {
	flags := flag.NewFlagSet("identity-mail-reconcile", flag.ContinueOnError)
	flags.SetOutput(diagnostics)
	tenant := flags.String("tenant-id", "", "Exact tenant identifier")
	database := flags.String("database", "", "Expected database name")
	after := flags.String("after", "", "Continuation delivery UUID")
	limit := flags.Int("limit", 100, "Maximum rows (1-500)")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 || *tenant == "" || strings.TrimSpace(*tenant) != *tenant || *database == "" {
		return errors.New("exact --tenant-id and --database are required")
	}
	if os.Getenv("DATABASE_URL") == "" {
		return errors.New("DATABASE_URL is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db, closeDB, err := s.database.Business(ctx)
	if err != nil {
		return errors.New("could not connect to configured database")
	}
	defer closeDB()

	var current string
	if err := db.WithContext(ctx).Raw("SELECT current_database()").Scan(&current).Error; err != nil {
		return errors.New("could not verify database name")
	}
	if current != *database {
		return errors.New("configured database does not match --database")
	}

	report, err := integration.NewMailReconciliationClient(db, identity.NewInvitationDeliveryReader(db)).Search(ctx, mail.SearchRequest{
		TenantID: *tenant, After: *after, Limit: *limit,
	})
	if err != nil {
		return errors.New("reconciliation query failed; verify scope and migrations")
	}

	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}
