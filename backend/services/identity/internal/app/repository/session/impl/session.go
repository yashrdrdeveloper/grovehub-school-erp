package impl
import ("context";"errors";"time";"gorm.io/gorm";tx "grovehub/libs/core/database/transaction";"grovehub/services/identity/internal/app/entity";contract "grovehub/services/identity/internal/app/repository/session")
type Repository struct{db *gorm.DB}
func New(db *gorm.DB)contract.Repository{return &Repository{db:db}}
func (r *Repository) Create(ctx context.Context,s entity.Session)error{return tx.Database(ctx,r.db).Create(&s).Error}
func (r *Repository) FindAccess(ctx context.Context,hash string,now time.Time)(entity.Session,error){
 var s entity.Session
 err:=tx.Database(ctx,r.db).Where("access_hash = ? AND revoked_at IS NULL AND access_expires_at > ?",hash,now).Take(&s).Error
 if errors.Is(err,gorm.ErrRecordNotFound){return s,contract.ErrInvalidSession};return s,err
}
func (r *Repository) FindRefresh(ctx context.Context,hash string,now time.Time)(entity.Session,error){
 var s entity.Session
 err:=tx.Database(ctx,r.db).Where("refresh_hash = ? AND revoked_at IS NULL AND refresh_expires_at > ?",hash,now).Take(&s).Error
 if errors.Is(err,gorm.ErrRecordNotFound){return s,contract.ErrInvalidSession};return s,err
}
func (r *Repository) Rotate(ctx context.Context,s entity.Session,oldHash string)error{
 result:=tx.Database(ctx,r.db).Model(&entity.Session{}).Where("session_id = ? AND refresh_hash = ? AND revoked_at IS NULL AND refresh_expires_at > ?",s.SessionID,oldHash,time.Now().UTC()).
 Updates(map[string]any{"access_hash":s.AccessHash,"refresh_hash":s.RefreshHash,"access_expires_at":s.AccessExpiresAt,"updated_at":s.UpdatedAt})
 if result.Error!=nil{return result.Error};if result.RowsAffected!=1{return contract.ErrInvalidSession};return nil
}
func (r *Repository) Reauthenticate(ctx context.Context,id string,when time.Time)error{
 res:=tx.Database(ctx,r.db).Model(&entity.Session{}).Where("session_id = ? AND revoked_at IS NULL",id).Update("reauthenticated_at",when)
 if res.Error!=nil{return res.Error};if res.RowsAffected!=1{return contract.ErrInvalidSession};return nil
}
func (r *Repository) List(ctx context.Context,accountID string)([]entity.Session,error){var s []entity.Session;err:=tx.Database(ctx,r.db).Where("account_id = ? AND revoked_at IS NULL AND refresh_expires_at > ?",accountID,time.Now().UTC()).Order("created_at DESC").Limit(100).Find(&s).Error;return s,err}
func (r *Repository) Revoke(ctx context.Context,accountID,id string)error{
 res:=tx.Database(ctx,r.db).Model(&entity.Session{}).Where("account_id = ? AND session_id = ? AND revoked_at IS NULL",accountID,id).Update("revoked_at",time.Now().UTC())
 if res.Error!=nil{return res.Error};if res.RowsAffected!=1{return contract.ErrInvalidSession};return nil
}
func (r *Repository) RevokeAll(ctx context.Context,accountID string)(int64,error){
 res:=tx.Database(ctx,r.db).Model(&entity.Session{}).Where("account_id = ? AND revoked_at IS NULL",accountID).Update("revoked_at",time.Now().UTC())
 return res.RowsAffected,res.Error
}
