package entity

import "time"

type DKIMKey struct {
	Domain        string    `gorm:"column:domain;primaryKey"`
	Selector      string    `gorm:"column:selector"`
	PrivateKeyPEM string    `gorm:"column:private_key_pem"`
	PublicKeyPEM  string    `gorm:"column:public_key_pem"`
	CreatedAt     time.Time `gorm:"column:created_at"`
}

func (DKIMKey) TableName() string { return "infrastructure.dkim_keys" }
