package impl

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"grovehub/launcher/internal/src/operations/command/registry"
	"grovehub/libs/platform/command"
)

// @Service
type dispatcher struct {
	commands map[string]command.Command
}

func New(commands map[string]command.Command) registry.Registry {
	return &dispatcher{commands: commands}
}

func (r *dispatcher) Run(args []string, output, diagnostics io.Writer) int {
	name := "serve"
	if len(args) > 0 {
		name, args = args[0], args[1:]
	}

	strategy, ok := r.commands[name]
	if !ok || strategy == nil {
		names := make([]string, 0, len(r.commands))
		for name := range r.commands {
			names = append(names, name)
		}
		sort.Strings(names)
		_, _ = fmt.Fprintf(diagnostics, "usage: groovehub-platform [%s]\n", strings.Join(names, "|"))
		return 2
	}

	return strategy.Run(args, output, diagnostics)
}
