package impl

import (
 "context"
 "time"
 "github.com/google/uuid"
 "gorm.io/gorm"
 "gorm.io/gorm/clause"
 "grovehub/services/integration/entity"
)
type WatchRepository struct{db *gorm.DB}
func NewWatchRepository(db *gorm.DB)*WatchRepository{return &WatchRepository{db:db}}
func(r *WatchRepository)SaveWatch(ctx context.Context,connectionID,historyID string,expires time.Time)error{
 entry:=entity.WebhookSubscription{
  ID:uuid.NewString(),ConnectionID:connectionID,Resource:"GMAIL_MAILBOX",
  ExternalID:historyID,Status:"ACTIVE",ExpiresAt:&expires,
 }
 return r.db.WithContext(ctx).Clauses(clause.OnConflict{
  Columns:[]clause.Column{{Name:"connection_id"},{Name:"resource"}},
  DoUpdates:clause.AssignmentColumns([]string{"external_id","status","expires_at"}),
 }).Create(&entry).Error
}
func(r *WatchRepository)GetGmailByMailbox(ctx context.Context,email string)([]entity.Connection,error){
 rows:=make([]entity.Connection,0)
 err:=r.db.WithContext(ctx).Where("provider = ? AND status = ? AND lower(external_id) = ?", "GMAIL",entity.StatusActive,email).Find(&rows).Error
 return rows,err
}
func(r *WatchRepository)ExpiringWatches(ctx context.Context,before time.Time)([]entity.Connection,error){
 rows:=make([]entity.Connection,0)
 err:=r.db.WithContext(ctx).Table("integration.connections AS c").
  Select("c.*").Joins("JOIN integration.webhook_subscriptions AS w ON w.connection_id = c.id").
  Where("c.provider = ? AND c.status = ? AND w.resource = ? AND w.status = ? AND w.expires_at < ?",
   "GMAIL",entity.StatusActive,"GMAIL_MAILBOX","ACTIVE",before).
  Limit(500).Scan(&rows).Error
 return rows,err
}
