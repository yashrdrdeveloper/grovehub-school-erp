package entity

import "time"

// MailReference is a provider-side deduplication ledger. Message content and
// user-facing threads are owned by Communication, not Integration.
type MailReference struct {
 ID string `gorm:"column:id;type:uuid;primaryKey"`
 ConnectionID string `gorm:"column:connection_id;type:uuid;not null;uniqueIndex:uq_mail_ref,priority:1"`
 ProviderMessageID string `gorm:"column:provider_message_id;not null;uniqueIndex:uq_mail_ref,priority:2"`
 ThreadID string `gorm:"column:thread_id;not null"`
 Direction string `gorm:"column:direction;not null"`
 HistoryID string `gorm:"column:history_id"`
 Labels string `gorm:"column:labels;not null"`
 Deleted bool `gorm:"column:deleted;not null"`
 UpdatedAt time.Time `gorm:"column:updated_at;not null"`
}
func (MailReference)TableName()string{return "integration.mail_references"}
