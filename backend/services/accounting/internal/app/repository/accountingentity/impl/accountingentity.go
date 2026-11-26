package impl

import (
 "context"
 "errors"

 "gorm.io/gorm"
 coretransaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
 "grovehub/services/accounting/internal/app/entity"
 repo "grovehub/services/accounting/internal/app/repository/accountingentity"
)

type Repository struct { db *gorm.DB }

var _ repo.Repository = (*Repository)(nil)

func New(db *gorm.DB) repo.Repository { return &Repository{db: db} }

func (r *Repository) Create(ctx context.Context, value entity.AccountingEntity) error {
 if value.TenantID == "" { return errors.New("tenant scope required") }
 return coretransaction.Database(ctx, r.db).Table("accounting.accounting_entities").Create(&value).Error
}

func (r *Repository) FindByID(ctx context.Context, tenantID, id string) (entity.AccountingEntity, error) {
 var value entity.AccountingEntity
 if tenantID == "" { return value, errors.New("tenant scope required") }
 err := coretransaction.Database(ctx, r.db).Table("accounting.accounting_entities").Where("tenant_id = ? AND entity_id = ?", tenantID, id).Take(&value).Error
 return value, err
}
