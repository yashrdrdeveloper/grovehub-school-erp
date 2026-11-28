// Package command defines the transport-independent operator boundary.
package command

import (
	"context"
	"io"
)

// Command executes isolated arguments and returns a process exit code.
type Command interface {
	Run(args []string, output, diagnostics io.Writer) int
}

// Database describes operator-owned connections without coupling the library to GORM.
type Database[D any] interface {
	Business(context.Context) (D, func(), error)
	Audit(context.Context, D) (D, func(), error)
	CurrentName(context.Context, D) (string, error)
}
