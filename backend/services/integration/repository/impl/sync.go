package impl

import (
 "context"
 "errors"
 "time"

 "github.com/google/uuid"
 "gorm.io/gorm"
 "gorm.io/gorm/clause"

 "grovehub/services/integration/entity"
)

type MailSyncRepository struct{db *gorm.DB}
func NewMailSyncRepository(db *gorm.DB)*MailSyncRepository{return &MailSyncRepository{db:db}}
func(r *MailSyncRepository)Cursor(ctx context.Context,connectionID string)(string,error){
 var state entity.SyncState
 err:=r.db.WithContext(ctx).Where("connection_id = ? AND resource = ?",connectionID,"MAIL").Take(&state).Error
 if errors.Is(err,gorm.ErrRecordNotFound){return "",nil}
 return state.Cursor,err
}
func(r *MailSyncRepository)SaveMessages(ctx context.Context,connectionID string,messages []entity.MailReference)error{
 if len(messages)==0{return nil}
 for i:=range messages {
  messages[i].ConnectionID=connectionID
  if messages[i].ID==""{messages[i].ID=uuid.NewString()}
  messages[i].UpdatedAt=time.Now().UTC()
 }
 return r.db.WithContext(ctx).Clauses(clause.OnConflict{
  Columns:[]clause.Column{{Name:"connection_id"},{Name:"provider_message_id"}},
  DoUpdates:clause.AssignmentColumns([]string{"thread_id","direction","history_id","labels","deleted","updated_at"}),
 }).CreateInBatches(&messages,100).Error
}
func(r *MailSyncRepository)MarkDeleted(ctx context.Context,connectionID string,ids []string)error{
 if len(ids)==0{return nil}
 return r.db.WithContext(ctx).Model(&entity.MailReference{}).
  Where("connection_id = ? AND provider_message_id IN ?",connectionID,ids).
  Updates(map[string]any{"deleted":true,"updated_at":time.Now().UTC()}).Error
}
func(r *MailSyncRepository)SaveCursor(ctx context.Context,connectionID,cursor string)error{
 state:=entity.SyncState{
 ID:uuid.NewString(),ConnectionID:connectionID,Resource:"MAIL",Cursor:cursor,
 }
 now:=time.Now().UTC()
 state.LastSyncedAt=&now
 return r.db.WithContext(ctx).Clauses(clause.OnConflict{
  Columns:[]clause.Column{{Name:"connection_id"},{Name:"resource"}},
  DoUpdates:clause.AssignmentColumns([]string{"cursor","last_synced_at"}),
 }).Create(&state).Error
}

func(r *MailSyncRepository)AcquireMailbox(ctx context.Context,connectionID string)(func(),bool,error){
 pool,err:=r.db.DB();if err!=nil{return nil,false,err}
 conn,err:=pool.Conn(ctx);if err!=nil{return nil,false,err}
 var acquired bool
 // Session advisory locks serialize sync across Go processes and instances.
 // Use dedicated sql.Conn so unlocking happens on the same PG session.
 err=conn.QueryRowContext(ctx,"SELECT pg_try_advisory_lock(hashtextextended($1,0))",connectionID).Scan(&acquired)
 if err!=nil{conn.Close();return nil,false,err}
 if !acquired{conn.Close();return nil,false,nil}
 release:=func(){
  cleanupCtx,cancel:=context.WithTimeout(context.Background(),5*time.Second)
  defer cancel()
  var unlocked bool
  _=conn.QueryRowContext(cleanupCtx,"SELECT pg_advisory_unlock(hashtextextended($1,0))",connectionID).Scan(&unlocked)
  _=conn.Close()
 }
 return release,true,nil
}
