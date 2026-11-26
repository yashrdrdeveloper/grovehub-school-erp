package registry

import (
	"context"
	"errors"

	"grovehub/services/identity/internal/app/constants/feature"
	"grovehub/services/identity/internal/app/constants/permission"
	registryrepo "grovehub/services/identity/internal/app/repository/registry"
)

var ErrInvalidBatch = errors.New("invalid capability registration batch")

type Request struct {
	Features []feature.Code
	Permissions []permission.Code
}

// Service is an internal operation. A future restricted workload adapter must
// authorize registration before calling it. It is not a browser API.
type Service interface {
	Register(ctx context.Context, request Request) (registryrepo.Result, error)
}
