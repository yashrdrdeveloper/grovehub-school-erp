package bootstrapadmin

import (
	"fmt"
	"io"
	"os"

	"grovehub/cmd/internal/src/command"
	admin "grovehub/services/identity/bootstrapadmin"
)

// @Strategy
type strategy struct{}

func New() command.Command {
	return strategy{}
}

func (strategy) Run(args []string, output, diagnostics io.Writer) int {
	if err := admin.RunArgs(args, os.Stdin, output, diagnostics); err != nil {
		_, _ = fmt.Fprintln(diagnostics, "bootstrap-admin:", err)
		return 1
	}

	return 0
}
