package entity

import "time"

type OAuthState struct {
 Hash string `gorm:"column:hash;primaryKey"`
 ConnectionID string `gorm:"column:connection_id;type:uuid;not null"`
 ExpiresAt time.Time `gorm:"column:expires_at;not null"`
}
func (OAuthState) TableName() string { return "integration.oauth_states" }

type OAuthCredential struct {
 ConnectionID string `gorm:"column:connection_id;type:uuid;primaryKey"`
 Ciphertext []byte `gorm:"column:ciphertext;not null"`
 UpdatedAt time.Time `gorm:"column:updated_at;not null"`
}
func (OAuthCredential) TableName() string { return "integration.oauth_credentials" }
