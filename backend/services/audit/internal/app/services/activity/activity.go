package activity

import (
	"context"
	"errors"
	dto "grovehub/services/audit/internal/app/dto/activity"
)

var ErrForbidden = errors.New("audit access denied")
var ErrInvalidSearch = errors.New("invalid audit search")

type Service interface {
	Search(ctx context.Context, query dto.Search) (dto.Page, error)
}
