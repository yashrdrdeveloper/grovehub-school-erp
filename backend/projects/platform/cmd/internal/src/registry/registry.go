package registry

import (
	"fmt"
	"io"

	"grovehub/cmd/internal/src/command"
	"grovehub/cmd/internal/src/constants"
)

// Registry dispatches an explicitly registered command.
type Registry interface {
	Run(args []string, output, diagnostics io.Writer) int
}

func New(commands map[string]command.Command) Registry {
	return &registry{commands: commands}
}

// @Service
type registry struct {
	commands map[string]command.Command
}

func (r *registry) Run(args []string, output, diagnostics io.Writer) int {
	name := constants.Serve
	if len(args) > 0 {
		name, args = args[0], args[1:]
	}

	strategy, ok := r.commands[name]
	if !ok || strategy == nil {
		_, _ = fmt.Fprintln(diagnostics, constants.Usage)
		return 2
	}

	return strategy.Run(args, output, diagnostics)
}
