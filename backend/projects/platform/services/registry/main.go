// Package registry exposes the project-independent Registry service bootstrap.
package registry

import (
	"grovehub/services/registry/internal/app/bootstrap"
	bootstrapimpl "grovehub/services/registry/internal/app/bootstrap/impl"
)

type Service = bootstrap.Service

func Bootstrap() Service {
	return bootstrapimpl.NewServiceBootstrap()
}
