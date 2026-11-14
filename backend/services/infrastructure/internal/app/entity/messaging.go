package entity

import (
    "encoding/json"
    "time"
)

type MessagingStatus string

const (
    MessagingPending MessagingStatus = "PENDING"
    MessagingSending MessagingStatus = "SENDING"
    MessagingSent MessagingStatus = "SENT"
)

// Messaging is the durable event delivery state, not a chat message.
// @Entity infrastructure.messaging
type Messaging struct {
    EventID string `gorm:"column:event_id;primaryKey;type:text"`
    Topic string `gorm:"column:topic;not null"`
    Payload []byte `gorm:"column:payload;not null"`
    Metadata json.RawMessage `gorm:"column:metadata;type:jsonb;not null"`
    Status MessagingStatus `gorm:"column:status;not null"`
    Attempts int `gorm:"column:attempts;not null"`
    NextAttemptAt time.Time `gorm:"column:next_attempt_at;not null"`
    LeaseUntil *time.Time `gorm:"column:lease_until"`
    CreatedAt time.Time `gorm:"column:created_at;not null"`
    PublishedAt *time.Time `gorm:"column:published_at"`
}
