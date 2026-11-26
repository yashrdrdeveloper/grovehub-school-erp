package entity

import "time"

// @Entity accounting.accounting_books
type AccountingBook struct {
 BookID string `gorm:"column:book_id;primaryKey;type:uuid"`
 TenantID string `gorm:"column:tenant_id;not null"`
 EntityID string `gorm:"column:entity_id;type:uuid;not null"`
 Code string `gorm:"column:code;not null"`
 Name string `gorm:"column:name;not null"`
 BaseCurrency string `gorm:"column:base_currency;not null"`
 CurrencyScale int16 `gorm:"column:currency_scale;not null"`
 AccountingBasis string `gorm:"column:accounting_basis;not null"`
 CreatedAt time.Time `gorm:"column:created_at"`
}
