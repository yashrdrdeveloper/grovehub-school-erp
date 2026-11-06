package impl

import (
 "context"
 "errors"
 "time"

 integrationresponse "grovehub/services/integration/dto/response"
 "grovehub/services/integration/entity"
 "grovehub/services/integration/providers/gmail"
)

var ErrPushUnavailable=errors.New("Gmail Pub/Sub push is not configured")

func(s *ConnectionService)Watch(ctx context.Context,tenantID,id string)(integrationresponse.Watch,error){
 if s.gmail==nil||s.pushTopic==""{return integrationresponse.Watch{},ErrPushUnavailable}
 c,err:=s.repository.Get(ctx,tenantID,id)
 if err!=nil{return integrationresponse.Watch{},err}
 if c.Provider!="GMAIL"||c.Status!=entity.StatusActive{return integrationresponse.Watch{},ErrConnectionInactive}
 token,err:=s.validToken(ctx,id);if err!=nil{return integrationresponse.Watch{},err}
 historyID,expiration,err:=s.gmail.Watch(ctx,token,s.pushTopic)
 if err!=nil{return integrationresponse.Watch{},err}
 if err=s.watchRepository.SaveWatch(ctx,id,historyID,expiration);err!=nil{return integrationresponse.Watch{},err}
 return integrationresponse.Watch{ConnectionID:id,ExpiresAt:expiration},nil
}
func(s *ConnectionService)ProcessPush(ctx context.Context,notification gmail.PushNotification)error{
 if s.gmail==nil||s.pushTopic==""{return ErrPushUnavailable}
 matches,err:=s.watchRepository.GetGmailByMailbox(ctx,notification.EmailAddress)
 if err!=nil{return err}
 // A Gmail push carries mailbox and historyId only. Always read server-side
 // credentials for an already active connection and never trust a tenant id.
 for _,connection:=range matches{
  _,err=s.Sync(ctx,connection.TenantID,connection.ID)
  if err!=nil{return err}
 }
 return nil
}
func(s *ConnectionService)RenewExpiring(ctx context.Context)error{
 if s.gmail==nil||s.pushTopic==""{return nil}
 rows,err:=s.watchRepository.ExpiringWatches(ctx,time.Now().UTC().Add(24*time.Hour))
 if err!=nil{return err}
 for _,c:=range rows{
  if _,err:=s.Watch(ctx,c.TenantID,c.ID);err!=nil{return err}
 }
 return nil
}
