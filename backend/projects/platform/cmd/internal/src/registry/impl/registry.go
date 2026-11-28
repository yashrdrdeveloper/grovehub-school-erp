package impl

import (
	"fmt"
	"io"

	"grovehub/cmd/internal/src/command"
	"grovehub/cmd/internal/src/constants"
	"grovehub/cmd/internal/src/registry"
)

// @Service
type dispatcher struct {
	commands map[string]command.Command
}

func New(commands map[string]command.Command) registry.Registry {
	return &dispatcher{commands: commands}
}

func (r *dispatcher) Run(args []string, output, diagnostics io.Writer) int {
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
