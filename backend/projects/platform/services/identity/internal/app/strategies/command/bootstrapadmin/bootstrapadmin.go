package bootstrapadmin

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"grovehub/libs/platform/command"
	admin "grovehub/services/identity/bootstrapadmin"
)

type Runner func([]string, io.Reader, io.Writer, io.Writer) error

// @Strategy
type strategy struct {
	run Runner
}

func New() command.Command {
	return NewWithRunner(admin.RunArgs)
}

func NewWithRunner(run Runner) command.Command {
	return &strategy{run: run}
}

func (s *strategy) Run(args []string, output, diagnostics io.Writer) int {
	if err := s.run(args, os.Stdin, output, diagnostics); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}

		_, _ = fmt.Fprintln(diagnostics, "bootstrap-admin:", err)
		return 1
	}

	return 0
}
