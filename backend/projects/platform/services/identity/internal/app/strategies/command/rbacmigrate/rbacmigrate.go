package rbacmigrate

import (
	"fmt"
	"io"

	"grovehub/libs/platform/command"
	"grovehub/services/identity/rbacmigrationcli"
)

type Runner func([]string, io.Writer, io.Writer) error

// @Strategy
type strategy struct {
	run Runner
}

func New() command.Command {
	return NewWithRunner(rbacmigrationcli.Run)
}

func NewWithRunner(run Runner) command.Command {
	return &strategy{run: run}
}

func (s *strategy) Run(args []string, output, diagnostics io.Writer) int {
	if err := s.run(args, output, diagnostics); err != nil {
		_, _ = fmt.Fprintln(diagnostics, "identity-rbac-migrate:", err)
		return 1
	}

	return 0
}
