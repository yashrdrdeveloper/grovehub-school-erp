package impl
import ("context";"errors";"gorm.io/gorm";tx "grovehub/libs/core/database/transaction";"grovehub/services/identity/internal/app/entity";contract "grovehub/services/identity/internal/app/repository/credential")
type Repository struct{db *gorm.DB}
func New(db *gorm.DB)contract.Repository{return &Repository{db:db}}
func (r *Repository)Find(ctx context.Context,id string)(entity.PasswordCredential,error){
 var c entity.PasswordCredential
 err:=tx.Database(ctx,r.db).Where("account_id = ?",id).Take(&c).Error
 if errors.Is(err,gorm.ErrRecordNotFound){return c,contract.ErrNotFound};return c,err
}
func (r *Repository)Create(ctx context.Context,c entity.PasswordCredential)error{return tx.Database(ctx,r.db).Create(&c).Error}
func (r *Repository)Update(ctx context.Context,c entity.PasswordCredential)error{
 res:=tx.Database(ctx,r.db).Model(&entity.PasswordCredential{}).Where("account_id = ? AND version = ?",c.AccountID,c.Version-1).
 Updates(map[string]any{"password_hash":c.PasswordHash,"password_changed_at":c.PasswordChangedAt,"version":c.Version,"updated_at":c.UpdatedAt})
 if res.Error!=nil{return res.Error};if res.RowsAffected!=1{return errors.New("credential version conflict")};return nil
}
