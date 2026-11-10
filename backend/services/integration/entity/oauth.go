package entity

import "time"

// @Entity integration.oauth_states mode=table
type OAuthState struct {
	Hash         string    `gorm:"column:hash;primaryKey"`
	ConnectionID string    `gorm:"column:connection_id;type:uuid;not null"`
	ExpiresAt    time.Time `gorm:"column:expires_at;not null"`
}


// @Entity integration.oauth_credentials mode=table
type OAuthCredential struct {
	ConnectionID string    `gorm:"column:connection_id;type:uuid;primaryKey"`
	Ciphertext   []byte    `gorm:"column:ciphertext;not null"`
	UpdatedAt    time.Time `gorm:"column:updated_at;not null"`
}

