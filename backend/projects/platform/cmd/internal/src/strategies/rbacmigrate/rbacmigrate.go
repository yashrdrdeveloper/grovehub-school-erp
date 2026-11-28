package rbacmigrate

import (
	"fmt"
	"io"

	"grovehub/cmd/internal/src/command"
	"grovehub/services/identity/rbacmigrationcli"
)

// @Strategy
type strategy struct{}

func New() command.Command {
	return strategy{}
}

func (strategy) Run(args []string, output, diagnostics io.Writer) int {
	if err := rbacmigrationcli.Run(args, output, diagnostics); err != nil {
		_, _ = fmt.Fprintln(diagnostics, "identity-rbac-migrate:", err)
		return 1
	}

	return 0
}
