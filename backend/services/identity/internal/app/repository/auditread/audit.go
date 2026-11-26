package auditread

import (
	"context"

	dto "grovehub/services/identity/internal/app/dto/audit"
)

type Reader interface {
	Search(ctx context.Context, tenantID string, query dto.Search) ([]dto.Event, error)
}
