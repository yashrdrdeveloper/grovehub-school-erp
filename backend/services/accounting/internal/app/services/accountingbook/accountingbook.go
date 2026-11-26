package accountingbook

import (
 "context"
 "grovehub/services/accounting/internal/app/dto"
 "grovehub/services/accounting/internal/app/entity"
)

type Service interface {
 Create(ctx context.Context, tenantID string, input dto.CreateAccountingBook) (entity.AccountingBook, error)
 FindByID(ctx context.Context, tenantID, id string) (entity.AccountingBook, error)
}
