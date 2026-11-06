package impl

import (
 "context"
 "errors"
 "encoding/json"
 "github.com/google/uuid"
 coreevents "grovehub/libs/core/events"
 mailboxevents "grovehub/services/integration/events/mailbox"
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
 // No message ingestion or cursor updates occur here. Communication will
 // consume the event and call Integration provider capabilities.
 for _,connection:=range matches{
  if s.messaging==nil{return errors.New("integration messaging is unavailable")}
  eventID:=uuid.NewString()
  data:=mailboxevents.Changed{
   Provider:"GMAIL",ConnectionID:connection.ID,TenantID:connection.TenantID,
   OwnerGID:connection.OwnerGID,Mailbox:notification.EmailAddress,
   HistoryHint:notification.HistoryID,
  }
  envelope:=coreevents.NewEnvelope(ctx,coreevents.Metadata{
   EventID:eventID,EventType:mailboxevents.ChangedEventType,Version:1,
   AggregateID:connection.ID,OccurredAt:time.Now().UTC(),
  },data)
  payload,err:=json.Marshal(envelope)
  if err!=nil{return err}
  if err=s.messaging.Publish(ctx,mailboxevents.ChangedTopic,payload,map[string]string{
   "eventId":eventID,"eventType":mailboxevents.ChangedEventType,"source":"integration",
  });err!=nil{return err}
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
