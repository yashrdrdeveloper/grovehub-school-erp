package accountingbook

import (
 "context"
 "grovehub/services/accounting/internal/app/entity"
)

// @Repository entity=AccountingBook
type Repository interface {
 Create(ctx context.Context, value entity.AccountingBook) error
 FindByID(ctx context.Context, tenantID, id string) (entity.AccountingBook, error)
}
