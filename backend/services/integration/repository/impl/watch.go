package impl

import (
 "context"
 "strings"
 "time"

 "github.com/google/uuid"
 "gorm.io/gorm"
 "gorm.io/gorm/clause"

 "grovehub/services/integration/constants"
 "grovehub/services/integration/entity"
)

// WatchRepository uses GORM's model, association, and upsert APIs.
// There are no manual JOIN/SELECT/SCAN expressions in this repository.
type WatchRepository struct{db *gorm.DB}

func NewWatchRepository(db *gorm.DB)*WatchRepository{
 return &WatchRepository{db:db}
}

func(r *WatchRepository)SaveWatch(ctx context.Context,connectionID,historyID string,expires time.Time)error{
 entry:=entity.WebhookSubscription{
  ID:uuid.NewString(),
  ConnectionID:connectionID,
  Resource:constants.ResourceGmailMailbox,
  ExternalID:historyID,
  Status:constants.StatusActive,
  ExpiresAt:&expires,
 }
 return r.db.WithContext(ctx).Clauses(clause.OnConflict{
  Columns:[]clause.Column{{Name:"connection_id"},{Name:"resource"}},
  DoUpdates:clause.AssignmentColumns([]string{"external_id","status","expires_at"}),
 }).Create(&entry).Error
}

func(r *WatchRepository)GetGmailByMailbox(ctx context.Context,email string)([]entity.Connection,error){
 email=strings.ToLower(strings.TrimSpace(email))
 rows:=make([]entity.Connection,0)
 // Provider account addresses are normalized to lowercase during OAuth.
 criteria:=entity.Connection{Provider:constants.ProviderGmail,Status:entity.StatusActive,ExternalID:&email}
 err:=r.db.WithContext(ctx).Model(&entity.Connection{}).Where(&criteria).Find(&rows).Error
 return rows,err
}

func(r *WatchRepository)ExpiringWatches(ctx context.Context,before time.Time)([]entity.Connection,error){
 subscriptions:=make([]entity.WebhookSubscription,0)
 err:=r.db.WithContext(ctx).Model(&entity.WebhookSubscription{}).
  Where(&entity.WebhookSubscription{Resource:constants.ResourceGmailMailbox,Status:constants.StatusActive}).
  Where("expires_at < ?",before).
  Preload("Connection").
  Limit(500).
  Find(&subscriptions).Error
 if err!=nil{return nil,err}

 result:=make([]entity.Connection,0,len(subscriptions))
 for _,subscription:=range subscriptions{
  connection:=subscription.Connection
  if connection.ID!=""&&connection.Provider==constants.ProviderGmail&&connection.Status==entity.StatusActive{
   result=append(result,connection)
  }
 }
 return result,nil
}
