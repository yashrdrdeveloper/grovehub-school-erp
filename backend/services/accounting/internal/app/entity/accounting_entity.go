package entity

import "time"

// @Entity accounting.accounting_entities
type AccountingEntity struct {
 EntityID string `gorm:"column:entity_id;primaryKey;type:uuid"`
 TenantID string `gorm:"column:tenant_id;not null"`
 Code string `gorm:"column:code;not null"`
 Name string `gorm:"column:name;not null"`
 CreatedAt time.Time `gorm:"column:created_at"`
}
