package role

import (
	"context"

	accountdto "grovehub/services/identity/internal/app/dto/account"
)

// CatalogReader lists active tenant roles. It does not authorize assignments.
type CatalogReader interface {
	ListActive(ctx context.Context, tenantID string) ([]accountdto.Role, error)
}
