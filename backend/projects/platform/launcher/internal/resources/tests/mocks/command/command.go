package command

import "io"

// @Mock
type Command interface {
	Run(args []string, output, diagnostics io.Writer) int
}
