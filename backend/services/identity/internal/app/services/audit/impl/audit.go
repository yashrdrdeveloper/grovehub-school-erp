package impl

import (
	"context"

	"grovehub/services/identity/internal/app/constants/permission"
	dto "grovehub/services/identity/internal/app/dto/audit"
	"grovehub/services/identity/internal/app/repository/auditread"
	"grovehub/services/identity/internal/app/repository/authority"
	"grovehub/services/identity/internal/app/services/audit"
	"grovehub/services/identity/internal/app/services/authorization"
	"grovehub/services/identity/internal/app/services/security/principal"
	validator "grovehub/services/identity/internal/app/validators/audit"
)

type Service struct {
	reader auditread.Reader
	locker authority.Locker
	authorizer authorization.Service
	legacy bool
}

func New(reader auditread.Reader, locker authority.Locker, authorizer authorization.Service, legacy bool) audit.Service {
	return &Service{reader: reader, locker: locker, authorizer: authorizer, legacy: legacy}
}

func (service *Service) Search(ctx context.Context, query dto.Search) (dto.Page, error) {
	actor, err := principal.RequireTenant(ctx)
	if err != nil { return dto.Page{}, err }
	if service.legacy { return dto.Page{}, authorization.ErrForbidden }
	if service.reader == nil || service.locker == nil || service.authorizer == nil { return dto.Page{}, authorization.ErrUnavailable }
	var result dto.Page
	err = service.locker.Run(ctx, actor.TenantID, func(tx context.Context) error {
		if err := service.authorizer.Require(tx, permission.AuditView); err != nil { return err }
		normalized, err := validator.Normalize(query)
		if err != nil { return err }
		rows, err := service.reader.Search(tx, actor.TenantID, normalized)
		if err != nil { return err }
		result.Items = rows
		if len(rows) > normalized.Limit {
			result.Items = rows[:normalized.Limit]
			last := result.Items[len(result.Items)-1]
			result.NextCursor = &dto.Cursor{OccurredAt: last.OccurredAt, Source: last.Source, EventID: last.EventID}
		}
		return nil
	})
	if err != nil { return dto.Page{}, err }
	return result, nil
}
