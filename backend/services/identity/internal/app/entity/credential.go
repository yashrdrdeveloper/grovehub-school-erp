package entity
import "time"
// @Entity identity.password_credentials
type PasswordCredential struct {
 AccountID string `gorm:"column:account_id;primaryKey;type:uuid"`
 PasswordHash string `gorm:"column:password_hash;not null"`
 PasswordChangedAt time.Time `gorm:"column:password_changed_at"`
 Version int64 `gorm:"column:version"`
 CreatedAt time.Time `gorm:"column:created_at"`
 UpdatedAt time.Time `gorm:"column:updated_at"`
}
