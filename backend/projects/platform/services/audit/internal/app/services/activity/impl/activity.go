package impl

import (
	"context"

	coresearch "grovehub/libs/core/search"
	schoolgateway "grovehub/libs/platform/session"
	dto "grovehub/services/audit/internal/app/dto/activity"
	repo "grovehub/services/audit/internal/app/repository/activity"
	"grovehub/services/audit/internal/app/services/activity"
	"grovehub/services/audit/internal/app/services/activity/helpers"
	validator "grovehub/services/audit/internal/app/validators/activity"
)

type Service struct {
	repository repo.Repository
	access     schoolgateway.ProtectedRouteAuthorizer
}

func New(repository repo.Repository, access schoolgateway.ProtectedRouteAuthorizer) activity.Service {
	return &Service{repository: repository, access: access}
}

func (service *Service) Search(ctx context.Context, input dto.Search) (dto.Page, error) {
	actor, ok := schoolgateway.PrincipalFrom(ctx)
	if !ok || actor.TenantID == "" || service.access == nil || service.repository == nil || service.access.AuthorizeProtectedRoute(ctx, "POST /audit/protected/search") != nil {
		return dto.Page{}, activity.ErrForbidden
	}
	self := ""
	if input.Scope == "ME" {
		self = actor.AccountID
	} else if service.access.AuthorizeProtectedRoute(ctx, "AUDIT tenant search") != nil {
		return dto.Page{}, activity.ErrForbidden
	}
	query, err := validator.Normalize(input)
	if err != nil {
		return dto.Page{}, err
	}
	rows, err := service.repository.Search(ctx, actor.TenantID, self, query)
	if err != nil {
		return dto.Page{}, err
	}
	result := dto.Page{Items: rows.Items, PageInfo: coresearch.NewPageResponse(query, rows)}
	if rows.HasNext && len(rows.Items) > 0 {
		result.PageInfo.NextCursor, err = helpers.NextCursor(query, input.PageInfo, rows.Items[len(rows.Items)-1])
		if err != nil {
			return dto.Page{}, err
		}
	}
	return result, nil
}
