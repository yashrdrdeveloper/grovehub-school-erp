package impl
import ("context";"time";"gorm.io/gorm";tx "grovehub/libs/core/database/transaction";"grovehub/services/identity/internal/app/entity";contract "grovehub/services/identity/internal/app/repository/passwordtoken")
type Repository struct{db *gorm.DB}
func New(db *gorm.DB)contract.Repository{return &Repository{db:db}}
func (r *Repository) Issue(ctx context.Context,t entity.PasswordToken)error{return tx.Database(ctx,r.db).Create(&t).Error}
func (r *Repository) Consume(ctx context.Context,hash string,p entity.TokenPurpose,now time.Time)(string,error){
 var row struct{AccountID string `gorm:"column:account_id"`}
 res:=tx.Database(ctx,r.db).Raw("UPDATE identity.password_tokens SET used_at = ? WHERE token_hash = ? AND purpose = ? AND used_at IS NULL AND expires_at > ? RETURNING account_id",now,hash,p,now).Scan(&row)
 if res.Error!=nil{return "",res.Error};if row.AccountID==""{return "",contract.ErrInvalidToken};return row.AccountID,nil
}
func (r *Repository) Invalidate(ctx context.Context,accountID string,p entity.TokenPurpose)error{
 return tx.Database(ctx,r.db).Model(&entity.PasswordToken{}).Where("account_id = ? AND purpose = ? AND used_at IS NULL",accountID,p).Update("used_at",time.Now().UTC()).Error
}
