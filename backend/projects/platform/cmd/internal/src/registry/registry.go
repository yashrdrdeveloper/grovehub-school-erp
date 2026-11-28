package registry

import "io"

// Registry dispatches an explicitly registered command.
type Registry interface {
	Run(args []string, output, diagnostics io.Writer) int
}
