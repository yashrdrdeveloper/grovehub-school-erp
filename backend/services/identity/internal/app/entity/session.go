package entity
import "time"
// @Entity identity.sessions
type Session struct {
 SessionID string `gorm:"column:session_id;primaryKey;type:uuid" json:"sessionId"`
 AccountID string `gorm:"column:account_id;type:uuid" json:"-"`
 AccessHash string `gorm:"column:access_hash" json:"-"`
 RefreshHash string `gorm:"column:refresh_hash" json:"-"`
 AccessExpiresAt time.Time `gorm:"column:access_expires_at" json:"-"`
 RefreshExpiresAt time.Time `gorm:"column:refresh_expires_at" json:"expiresAt"`
 ReauthenticatedAt *time.Time `gorm:"column:reauthenticated_at" json:"-"`
 UserAgent string `gorm:"column:user_agent" json:"-"`
 LastActiveAt *time.Time `gorm:"column:last_active_at" json:"lastActiveAt,omitempty"`
 RevokedAt *time.Time `gorm:"column:revoked_at" json:"-"`
 CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
 UpdatedAt time.Time `gorm:"column:updated_at" json:"-"`
}
