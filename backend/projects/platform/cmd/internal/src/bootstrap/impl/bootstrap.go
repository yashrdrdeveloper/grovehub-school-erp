package impl

import (
	"grovehub/cmd/internal/src/bootstrap"
	"grovehub/cmd/internal/src/operations/execution"
	executionimpl "grovehub/cmd/internal/src/operations/execution/impl"
	"grovehub/launcher"
)

type builder struct{}

var _ bootstrap.Bootstrap = (*builder)(nil)

func New() bootstrap.Bootstrap {
	return &builder{}
}

func (b *builder) Build() execution.Executor {
	return executionimpl.New(launcher.Execute)
}
