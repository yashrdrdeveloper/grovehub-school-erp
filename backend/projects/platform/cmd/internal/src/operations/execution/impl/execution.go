package impl

import (
	"io"

	"grovehub/cmd/internal/src/operations/execution"
)

type Dispatch func([]string, io.Writer, io.Writer) int

// @Strategy
type executor struct {
	dispatch Dispatch
}

var _ execution.Executor = (*executor)(nil)

func New(dispatch Dispatch) execution.Executor {
	return &executor{dispatch: dispatch}
}

func (e *executor) Run(args []string, output, diagnostics io.Writer) int {
	return e.dispatch(args, output, diagnostics)
}
