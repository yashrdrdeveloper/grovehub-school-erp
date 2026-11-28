package activity

import (
	"context"
	coreerrors "grovehub/libs/core/errors"
	definitions "grovehub/services/audit/internal/app/constants/error"
	dto "grovehub/services/audit/internal/app/dto/activity"
)

var ErrForbidden = coreerrors.NewError(definitions.Forbidden)
var ErrInvalidSearch = coreerrors.NewError(definitions.InvalidRequest)

type Service interface {
	Search(ctx context.Context, query dto.Search) (dto.Page, error)
}
