package impl

import (
	"context"

	"grovehub/launcher/internal/src/bootstrap"
	"grovehub/launcher/internal/src/bootstrap/helpers"
	"grovehub/launcher/internal/src/server"
)

type runtime struct {
	ctx     context.Context
	server  server.Server
	cleanup *helpers.Cleanup
}

var _ bootstrap.Runtime = (*runtime)(nil)

func (r *runtime) Run() error {
	return r.server.Run(r.ctx)
}

func (r *runtime) Close() {
	r.cleanup.Close()
}
