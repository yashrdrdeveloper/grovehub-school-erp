package entity

import "time"

// PasswordSetupToken contains only the SHA-256 digest of a random setup token.
// The plaintext bearer secret must never be saved to this table or logs.
// @Entity identity.password_setup_tokens
type PasswordSetupToken struct {
    TokenHash string `gorm:"column:token_hash;primaryKey;type:char(64)"`
    AccountId AccountId `gorm:"column:account_id;type:uuid;not null"`
    ExpiresAt time.Time `gorm:"column:expires_at;not null"`
    UsedAt *time.Time `gorm:"column:used_at"`
    CreatedAt time.Time `gorm:"column:created_at;not null"`
}

func (PasswordSetupToken) TableName() string { return "identity.password_setup_tokens" }
