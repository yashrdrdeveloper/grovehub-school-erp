package accountingentity

import (
 "context"
 "grovehub/services/accounting/internal/app/dto"
 "grovehub/services/accounting/internal/app/entity"
)

type Service interface {
 Create(ctx context.Context, tenantID string, input dto.CreateAccountingEntity) (entity.AccountingEntity, error)
 FindByID(ctx context.Context, tenantID, id string) (entity.AccountingEntity, error)
}
