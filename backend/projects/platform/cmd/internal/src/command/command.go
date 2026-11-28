package command

import "io"

// Command is an operator strategy with isolated arguments and output streams.
type Command interface {
	Run(args []string, output, diagnostics io.Writer) int
}
