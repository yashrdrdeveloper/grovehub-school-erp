package audit

import (
	"context"
	"errors"

	dto "grovehub/services/identity/internal/app/dto/audit"
)

var ErrInvalidSearch = errors.New("invalid audit search")

type Service interface {
	Search(ctx context.Context, query dto.Search) (dto.Page, error)
}
