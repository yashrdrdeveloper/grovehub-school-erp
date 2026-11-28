package helpers

import (
	"flag"
	"io"

	bootstrap "grovehub/services/audit/internal/app/bootstrap/constants"
	"grovehub/services/audit/internal/app/strategies/command/auditstorage/constants"
	shared "grovehub/services/audit/internal/app/strategies/command/constants"
)

func Parse(args []string, diagnostics io.Writer) (string, error) {
	flags := flag.NewFlagSet(bootstrap.AuditStorage, flag.ContinueOnError)
	flags.SetOutput(diagnostics)
	database := flags.String(shared.DatabaseFlag, "", constants.DatabaseDescription)

	if err := flags.Parse(args); err != nil {
		return "", err
	}
	if flags.NArg() != 0 || *database == "" {
		return "", constants.ErrDatabaseRequired
	}

	return *database, nil
}
