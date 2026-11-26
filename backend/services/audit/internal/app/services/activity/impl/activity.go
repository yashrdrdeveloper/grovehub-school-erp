package impl

import (
	"context"

	"grovehub/libs/schoolgateway"
	dto "grovehub/services/audit/internal/app/dto/activity"
	repo "grovehub/services/audit/internal/app/repository/activity"
	"grovehub/services/audit/internal/app/services/activity"
	validator "grovehub/services/audit/internal/app/validators/activity"
)

type Service struct {
	repository repo.Repository
	access schoolgateway.ProtectedRouteAuthorizer
}

func New(repository repo.Repository, access schoolgateway.ProtectedRouteAuthorizer) activity.Service {
	return &Service{repository:repository, access:access}
}

func (service *Service) Search(ctx context.Context, query dto.Search) (dto.Page, error) {
	actor, ok := schoolgateway.PrincipalFrom(ctx)
	if !ok || actor.TenantID == "" || service.access == nil || service.access.AuthorizeProtectedRoute(ctx, "POST /audit/protected/search") != nil { return dto.Page{}, activity.ErrForbidden }
	query, err := validator.Normalize(query)
	if err != nil { return dto.Page{},err }
	self := ""
	if query.Scope == "ME" { self = actor.AccountID }
	rows, err := service.repository.Search(ctx, actor.TenantID, self, query)
	if err != nil { return dto.Page{},err }
	result := dto.Page{Items:rows}
	if len(rows) > query.Limit {
		result.Items = rows[:query.Limit]
		last := result.Items[len(result.Items)-1]
		result.NextCursor = &dto.Cursor{OccurredAt:last.OccurredAt, AuditID:last.AuditID}
	}
	return result,nil
}
