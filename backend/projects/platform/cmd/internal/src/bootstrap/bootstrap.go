// Package bootstrap defines the composition boundary for the executable.
package bootstrap

import "grovehub/cmd/internal/src/operations/execution"

type Bootstrap interface {
	Build() execution.Executor
}
