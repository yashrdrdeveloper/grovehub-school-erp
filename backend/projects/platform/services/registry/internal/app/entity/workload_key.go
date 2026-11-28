package entity

import "time"

// @Entity registry.workload_keys
type WorkloadKey struct {
	KeyID      string     `gorm:"column:key_id;primaryKey"`
	WorkloadID string     `gorm:"column:workload_id;type:uuid;not null"`
	PublicKey  []byte     `gorm:"column:public_key;type:bytea;not null"`
	ValidFrom  time.Time  `gorm:"column:valid_from;not null"`
	ValidUntil time.Time  `gorm:"column:valid_until;not null"`
	RevokedAt  *time.Time `gorm:"column:revoked_at"`
	CreatedAt  time.Time  `gorm:"column:created_at;not null"`
}
