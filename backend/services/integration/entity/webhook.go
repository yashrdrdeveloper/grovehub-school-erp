package entity

import "time"

// WebhookSubscription stores provider connectivity, never mailbox messages.
type WebhookSubscription struct {
	ID           string     `gorm:"column:id;type:uuid;primaryKey"`
	ConnectionID string     `gorm:"column:connection_id;type:uuid;not null"`
	Resource     string     `gorm:"column:resource;not null"`
	ExternalID   string     `gorm:"column:external_id;not null"`
	Status       string     `gorm:"column:status;not null"`
	ExpiresAt    *time.Time `gorm:"column:expires_at"`
	Connection   Connection `gorm:"foreignKey:ConnectionID;references:ID"`
}

func (WebhookSubscription) TableName() string { return "integration.webhook_subscriptions" }
