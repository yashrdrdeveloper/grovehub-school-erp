package entity
import "time"
type TokenPurpose string
const ( AccountSetup TokenPurpose="ACCOUNT_SETUP"; PasswordReset TokenPurpose="PASSWORD_RESET"; SetupGrant TokenPurpose="SETUP_GRANT"; ResetGrant TokenPurpose="RESET_GRANT" )
// @Entity identity.password_tokens
type PasswordToken struct {
 TokenHash string `gorm:"column:token_hash;primaryKey;type:char(64)"`
 AccountID string `gorm:"column:account_id;type:uuid"`
 Purpose TokenPurpose `gorm:"column:purpose"`
 ExpiresAt time.Time `gorm:"column:expires_at"`
 UsedAt *time.Time `gorm:"column:used_at"`
 CreatedAt time.Time `gorm:"column:created_at"`
}
