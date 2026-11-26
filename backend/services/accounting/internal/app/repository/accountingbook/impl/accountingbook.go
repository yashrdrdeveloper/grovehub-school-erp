package impl

import (
 "context"
 "errors"

 "gorm.io/gorm"
 coretransaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
 "grovehub/services/accounting/internal/app/entity"
 repo "grovehub/services/accounting/internal/app/repository/accountingbook"
)

type Repository struct { db *gorm.DB }

var _ repo.Repository = (*Repository)(nil)

func New(db *gorm.DB) repo.Repository { return &Repository{db: db} }

func (r *Repository) Create(ctx context.Context, value entity.AccountingBook) error {
 if value.TenantID == "" { return errors.New("tenant scope required") }
 return coretransaction.Database(ctx, r.db).Table("accounting.accounting_books").Create(&value).Error
}

func (r *Repository) FindByID(ctx context.Context, tenantID, id string) (entity.AccountingBook, error) {
 var value entity.AccountingBook
 if tenantID == "" { return value, errors.New("tenant scope required") }
 err := coretransaction.Database(ctx, r.db).Table("accounting.accounting_books").Where("tenant_id = ? AND book_id = ?", tenantID, id).Take(&value).Error
 return value, err
}
