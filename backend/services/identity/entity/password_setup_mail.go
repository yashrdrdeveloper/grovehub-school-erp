package entity

import "time"

// PasswordSetupMail is an encrypted durable invitation delivery job.
// Ciphertext includes a one-time URL; encryption key is never persisted here.
// @Entity identity.password_setup_mail
type PasswordSetupMail struct {
    AccountId AccountId `gorm:"column:account_id;type:uuid;primaryKey"`
    Email string `gorm:"column:email;not null"`
    Ciphertext []byte `gorm:"column:ciphertext;not null"`
    Nonce []byte `gorm:"column:nonce;not null"`
    Attempts int `gorm:"column:attempts;not null"`
    NextAttemptAt time.Time `gorm:"column:next_attempt_at;not null"`
    LeaseUntil *time.Time `gorm:"column:lease_until"`
    SentAt *time.Time `gorm:"column:sent_at"`
    CreatedAt time.Time `gorm:"column:created_at;not null"`
}

func (PasswordSetupMail) TableName() string { return "identity.password_setup_mail" }
