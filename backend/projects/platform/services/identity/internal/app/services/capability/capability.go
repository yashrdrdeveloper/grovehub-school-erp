package capability

import (
	"context"

	dto "grovehub/services/identity/internal/app/dto/capability"
)

type Service interface {
	Discover(ctx context.Context) (dto.Discovery, error)
	Catalog(ctx context.Context) (dto.Catalog, error)
}
