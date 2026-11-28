package auditreplay

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
	"grovehub/services/audit"
	"grovehub/services/infrastructure"
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
		_, _ = fmt.Fprintln(diagnostics, "audit-replay:", err)
		return 1
	}

	return 0
}

func (s *strategy) run(args []string, output, diagnostics io.Writer) error {
	flags := flag.NewFlagSet("audit-replay", flag.ContinueOnError)
	flags.SetOutput(diagnostics)
	tenant := flags.String("tenant-id", "", "Exact tenant identifier")
	event := flags.String("event-id", "", "Exact original published event UUID")
	database := flags.String("database", "", "Expected business database name")
	auditDatabase := flags.String("audit-database", "", "Expected Audit database name; required with AUDIT_DATABASE_URL")
	operator := flags.String("operator", "", "Operator label (required for apply)")
	apply := flags.Bool("apply", false, "Persist the original audit event and replay receipt")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 || *tenant == "" || *event == "" || *database == "" || (*apply && *operator == "") {
		return errors.New("tenant-id, event-id, database and operator for apply are required")
	}
	if strings.TrimSpace(os.Getenv("AUDIT_DATABASE_URL")) != "" && *auditDatabase == "" {
		return errors.New("audit-database is required with AUDIT_DATABASE_URL")
	}
	if *auditDatabase == "" {
		*auditDatabase = *database
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db, closeDB, err := s.database.Business(ctx)
	if err != nil {
		return errors.New("could not connect to business database")
	}
	defer closeDB()

	var current string
	if err := db.WithContext(ctx).Raw("SELECT current_database()").Scan(&current).Error; err != nil || current != *database {
		return errors.New("configured business database does not match --database")
	}

	auditDB, closeAudit, err := s.database.Audit(ctx, db)
	if err != nil {
		return err
	}
	defer closeAudit()

	if err := auditDB.WithContext(ctx).Raw("SELECT current_database()").Scan(&current).Error; err != nil || current != *auditDatabase {
		return errors.New("configured Audit database does not match --audit-database")
	}

	result, err := audit.Replay(ctx, auditDB, infrastructure.NewReplaySource(db), audit.ReplayRequest{
		TenantID: *tenant, EventID: *event, OperatorLabel: *operator, Apply: *apply,
	})
	if err != nil {
		return errors.New("replay failed; verify retained published event, tenant, schema and identity conflicts")
	}

	return json.NewEncoder(output).Encode(result)
}
