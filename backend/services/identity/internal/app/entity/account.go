package entity
import "time"
type AccountStatus string
const ( PendingSetup AccountStatus="PENDING_SETUP"; Active AccountStatus="ACTIVE"; Suspended AccountStatus="SUSPENDED"; Deactivated AccountStatus="DEACTIVATED" )
// @Entity identity.accounts
type Account struct {
 AccountID string `gorm:"column:account_id;primaryKey;type:uuid" json:"accountId"`
 PersonGID string `gorm:"column:person_gid;type:uuid;not null" json:"personGid"`
 TenantID *string `gorm:"column:tenant_id" json:"-"`
 Username *string `gorm:"column:username" json:"username"`
 Email string `gorm:"column:email;not null" json:"email"`
 Phone *string `gorm:"column:phone" json:"-"`
 Status AccountStatus `gorm:"column:status;not null" json:"status"`
 Version int64 `gorm:"column:version;not null" json:"version"`
 SearchDetails string `gorm:"column:search_details;not null" json:"-"`
 RegistrationIdempotencyKey *string `gorm:"column:registration_idempotency_key" json:"-"`
 RegistrationRequestHash *string `gorm:"column:registration_request_hash" json:"-"`
 CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
 UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}
