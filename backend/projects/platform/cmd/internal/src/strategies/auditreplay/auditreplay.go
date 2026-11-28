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

	"gorm.io/gorm"

	"grovehub/cmd/internal/src/command"
	"grovehub/cmd/internal/src/operations/database"
	"grovehub/services/audit"
	"grovehub/services/infrastructure"
)

type ReplayOperation func(context.Context, *gorm.DB, *gorm.DB, audit.ReplayRequest) (audit.ReplayResult, error)

// @Strategy
type strategy struct {
	database database.Database
	replay   ReplayOperation
}

func New(db database.Database) command.Command {
	return NewWithReplay(db, func(ctx context.Context, auditDB, businessDB *gorm.DB, request audit.ReplayRequest) (audit.ReplayResult, error) {
		return audit.Replay(ctx, auditDB, infrastructure.NewReplaySource(businessDB), request)
	})
}

func NewWithReplay(db database.Database, replay ReplayOperation) command.Command {
	return &strategy{database: db, replay: replay}
}

func (s *strategy) Run(args []string, output, diagnostics io.Writer) int {
	if err := s.run(args, output, diagnostics); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}

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

	current, err := s.database.CurrentName(ctx, db)
	if err != nil || current != *database {
		return errors.New("configured business database does not match --database")
	}

	auditDB, closeAudit, err := s.database.Audit(ctx, db)
	if err != nil {
		return err
	}
	defer closeAudit()

	current, err = s.database.CurrentName(ctx, auditDB)
	if err != nil || current != *auditDatabase {
		return errors.New("configured Audit database does not match --audit-database")
	}

	result, err := s.replay(ctx, auditDB, db, audit.ReplayRequest{
		TenantID: *tenant, EventID: *event, OperatorLabel: *operator, Apply: *apply,
	})
	if err != nil {
		return errors.New("replay failed; verify retained published event, tenant, schema and identity conflicts")
	}

	return json.NewEncoder(output).Encode(result)
}
