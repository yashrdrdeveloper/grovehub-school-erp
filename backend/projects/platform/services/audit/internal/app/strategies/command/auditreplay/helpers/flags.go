package helpers

import (
	"flag"
	"io"
	"os"
	"strings"

	bootstrap "grovehub/services/audit/internal/app/bootstrap/constants"
	"grovehub/services/audit/internal/app/strategies/command/auditreplay/constants"
	shared "grovehub/services/audit/internal/app/strategies/command/constants"
	"grovehub/services/audit/internal/app/services/replay"
)

type Options struct {
	Database      string
	AuditDatabase string
	Request       replay.Request
}

func Parse(args []string, diagnostics io.Writer) (Options, error) {
	var options Options
	flags := flag.NewFlagSet(bootstrap.AuditReplay, flag.ContinueOnError)
	flags.SetOutput(diagnostics)
	flags.StringVar(&options.Request.TenantID, constants.TenantFlag, "", constants.TenantDescription)
	flags.StringVar(&options.Request.EventID, constants.EventFlag, "", constants.EventDescription)
	flags.StringVar(&options.Database, shared.DatabaseFlag, "", constants.DatabaseDescription)
	flags.StringVar(&options.AuditDatabase, constants.AuditDatabaseFlag, "", constants.AuditDatabaseDescription)
	flags.StringVar(&options.Request.OperatorLabel, constants.OperatorFlag, "", constants.OperatorDescription)
	flags.BoolVar(&options.Request.Apply, constants.ApplyFlag, false, constants.ApplyDescription)

	if err := flags.Parse(args); err != nil {
		return Options{}, err
	}

	if flags.NArg() != 0 || options.Request.TenantID == "" || options.Request.EventID == "" ||
		options.Database == "" || (options.Request.Apply && options.Request.OperatorLabel == "") {
		return Options{}, constants.ErrRequired
	}

	if strings.TrimSpace(os.Getenv(shared.DatabaseURL)) != "" && options.AuditDatabase == "" {
		return Options{}, constants.ErrAuditDatabaseRequired
	}
	if options.AuditDatabase == "" {
		options.AuditDatabase = options.Database
	}

	return options, nil
}
