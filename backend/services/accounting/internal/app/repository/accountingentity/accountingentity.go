package accountingentity

import (
 "context"
 "grovehub/services/accounting/internal/app/entity"
)

// @Repository entity=AccountingEntity
type Repository interface {
 Create(ctx context.Context, value entity.AccountingEntity) error
 FindByID(ctx context.Context, tenantID, id string) (entity.AccountingEntity, error)
}
