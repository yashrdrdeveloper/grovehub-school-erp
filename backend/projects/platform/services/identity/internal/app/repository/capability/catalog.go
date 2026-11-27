package capability

import (
	"context"

	"grovehub/services/identity/internal/app/constants/feature"
	"grovehub/services/identity/internal/app/constants/permission"
)

type Entry struct {
	Code permission.Code
	Feature feature.Code
}

type Reader interface {
	Available(ctx context.Context, tenantID string) ([]Entry, error)
}
