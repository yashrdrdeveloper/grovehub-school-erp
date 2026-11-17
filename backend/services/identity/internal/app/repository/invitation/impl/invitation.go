package impl
import ("context";"time";"gorm.io/gorm";tx "grovehub/libs/core/database/transaction";"grovehub/services/identity/internal/app/entity";contract "grovehub/services/identity/internal/app/repository/invitation")
type Repository struct{db *gorm.DB}
func New(db *gorm.DB)contract.Repository{return &Repository{db:db}}
func (r *Repository)Queue(ctx context.Context,m entity.InvitationMail)error{return tx.Database(ctx,r.db).Create(&m).Error}
func (r *Repository)Claim(ctx context.Context,now time.Time)(entity.InvitationMail,bool,error){
 var m entity.InvitationMail
 res:=tx.Database(ctx,r.db).Raw(`UPDATE identity.invitation_mail SET lease_until = ?, attempts = attempts + 1 WHERE id = (
 SELECT id FROM identity.invitation_mail WHERE sent_at IS NULL AND next_attempt_at <= ?
 AND (lease_until IS NULL OR lease_until < ?) ORDER BY next_attempt_at FOR UPDATE SKIP LOCKED LIMIT 1
 ) RETURNING *`,now.Add(45*time.Second),now,now).Scan(&m)
 return m,m.ID!="",res.Error
}
func (r *Repository)Delivered(ctx context.Context,id string)error{return tx.Database(ctx,r.db).Model(&entity.InvitationMail{}).Where("id = ?",id).Updates(map[string]any{"sent_at":time.Now().UTC(),"lease_until":nil}).Error}
func (r *Repository)Retry(ctx context.Context,id string,next time.Time)error{return tx.Database(ctx,r.db).Model(&entity.InvitationMail{}).Where("id = ?",id).Updates(map[string]any{"next_attempt_at":next,"lease_until":nil}).Error}
