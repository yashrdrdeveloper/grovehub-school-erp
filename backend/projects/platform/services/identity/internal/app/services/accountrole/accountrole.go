package accountrole

import (
	"context"

	dto "grovehub/services/identity/internal/app/dto/accountrole"
	repo "grovehub/services/identity/internal/app/repository/accountrole"
)

type Service interface {
	Get(ctx context.Context, accountID string) (repo.Assignment, error)
	Preview(ctx context.Context, accountID string, request dto.Replace) (dto.Preview, error)
}
