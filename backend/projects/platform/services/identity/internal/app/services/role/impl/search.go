package impl

import (
	"context"

	coresearch "grovehub/libs/core/search"
	"grovehub/services/identity/internal/app/constants/permission"
	roledto "grovehub/services/identity/internal/app/dto/role"
	searchoperation "grovehub/services/identity/internal/app/services/role/operations/search"
	rolevalidator "grovehub/services/identity/internal/app/validators/role"
)

func (service *Service) Search(ctx context.Context, request coresearch.Request) (coresearch.Result[roledto.Detail], error) {
	actor, err := service.authorize(ctx, permission.RoleView)
	if err != nil {
		return coresearch.Result[roledto.Detail]{}, err
	}
	query, err := rolevalidator.NormalizeSearch(request)
	if err != nil {
		return coresearch.Result[roledto.Detail]{}, err
	}
	rows, err := service.repository.Search(ctx, actor.TenantID, query)
	if err != nil {
		return coresearch.Result[roledto.Detail]{}, err
	}
	return searchoperation.Hydrate(query, rows)
}
