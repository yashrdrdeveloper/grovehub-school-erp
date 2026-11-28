package entity

import "time"

// @Entity identity.invitations
type Invitation struct {
	ID            string     `gorm:"column:id;primaryKey;type:uuid"`
	AccountID     string     `gorm:"column:account_id;type:uuid"`
	Email         string     `gorm:"column:email"`
	Purpose       string     `gorm:"column:purpose"`
	Ciphertext    []byte     `gorm:"column:ciphertext"`
	Nonce         []byte     `gorm:"column:nonce"`
	Attempts      int        `gorm:"column:attempts"`
	NextAttemptAt time.Time  `gorm:"column:next_attempt_at"`
	LeaseUntil    *time.Time `gorm:"column:lease_until"`
	SentAt        *time.Time `gorm:"column:sent_at"`
	SuppressedAt  *time.Time `gorm:"column:suppressed_at"`
	SuppressionReason *string `gorm:"column:suppression_reason"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
}
