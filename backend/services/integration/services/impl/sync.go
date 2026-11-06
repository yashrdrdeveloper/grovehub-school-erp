package impl

import (
 "context"
 "encoding/json"
 "errors"
 "fmt"
 "strings"

 integrationresponse "grovehub/services/integration/dto/response"
 "grovehub/services/integration/entity"
 "grovehub/services/integration/providers/gmail"
)

const maxSyncPages=100

// ErrSyncTooLarge means the initial scan was deliberately not acknowledged.
// It can be retried; upsert makes replay idempotent.
var ErrSyncTooLarge=errors.New("mailbox sync requires more pages than allowed per run")

func(s *ConnectionService)Sync(ctx context.Context,tenantID,id string)(integrationresponse.SyncResult,error){
 if s.gmail==nil{return integrationresponse.SyncResult{},ErrOAuthUnavailable}
 c,err:=s.repository.Get(ctx,tenantID,id)
 if err!=nil{return integrationresponse.SyncResult{},err}
 if c.Provider!="GMAIL"||c.Status!=entity.StatusActive{return integrationresponse.SyncResult{},ErrConnectionInactive}
 token,err:=s.validToken(ctx,id)
 if err!=nil{return integrationresponse.SyncResult{},err}
 cursor,err:=s.syncRepository.Cursor(ctx,id)
 if err!=nil{return integrationresponse.SyncResult{},err}
 if cursor=="" {
  return s.fullSync(ctx,id,token)
 }
 result,err:=s.incrementalSync(ctx,id,token,cursor)
 if errors.Is(err,gmail.ErrExpiredHistory){return s.fullSync(ctx,id,token)}
 return result,err
}
func(s *ConnectionService)fullSync(ctx context.Context,id string,token gmail.Token)(integrationresponse.SyncResult,error){
 // Take a baseline before listing, so updates during initial scan are picked
 // up on the next incremental call rather than silently skipped.
 baseline,err:=s.gmail.HistoryID(ctx,token)
 if err!=nil{return integrationresponse.SyncResult{},err}
 result:=integrationresponse.SyncResult{ConnectionID:id,Mode:"FULL"}
 page:=""
 for i:=0;i<maxSyncPages;i++{
  batch,err:=s.gmail.Messages(ctx,token,page)
  if err!=nil{return integrationresponse.SyncResult{},err}
  references:=make([]entity.MailReference,0,len(batch.Messages))
  for _,m:=range batch.Messages{
   detail,err:=s.gmail.MessageMetadata(ctx,token,m.ID)
   if err!=nil{return integrationresponse.SyncResult{},err}
   ref,err:=mailReference(detail)
   if err!=nil{return integrationresponse.SyncResult{},err}
   references=append(references,ref)
  }
  if err:=s.syncRepository.SaveMessages(ctx,id,references);err!=nil{return integrationresponse.SyncResult{},err}
  result.Processed+=len(references)
  page=batch.NextPageToken
  if page==""{
   if err:=s.syncRepository.SaveCursor(ctx,id,baseline);err!=nil{return integrationresponse.SyncResult{},err}
   result.Cursor=baseline
   return result,nil
  }
 }
 return integrationresponse.SyncResult{},ErrSyncTooLarge
}
func(s *ConnectionService)incrementalSync(ctx context.Context,id string,token gmail.Token,cursor string)(integrationresponse.SyncResult,error){
 result:=integrationresponse.SyncResult{ConnectionID:id,Mode:"INCREMENTAL"}
 page:=""
 for i:=0;i<maxSyncPages;i++{
  changes,err:=s.gmail.History(ctx,token,cursor,page)
  if err!=nil{return integrationresponse.SyncResult{},err}
  // A message may appear in several history records. Updating the same key
  // is harmless; only providerMessageId+connectionId identifies one item.
  upserts:=map[string]struct{}{}
  deletes:=map[string]struct{}{}
  for _,h:=range changes.History{
   for _,item:=range h.MessagesAdded{if item.Message.ID!=""{upserts[item.Message.ID]=struct{}{}}}
   for _,item:=range h.LabelsAdded{if item.Message.ID!=""{upserts[item.Message.ID]=struct{}{}}}
   for _,item:=range h.LabelsRemoved{if item.Message.ID!=""{upserts[item.Message.ID]=struct{}{}}}
   for _,item:=range h.MessagesDeleted{if item.Message.ID!=""{deletes[item.Message.ID]=struct{}{}}}
  }
  references:=make([]entity.MailReference,0,len(upserts))
  for messageID:=range upserts{
   if _,removed:=deletes[messageID];removed{continue}
   detail,err:=s.gmail.MessageMetadata(ctx,token,messageID)
   if err!=nil{
    // A concurrent delete can race metadata retrieval; it will be reconciled
    // on the next history page or sync. Do not advance cursor on an error.
    return integrationresponse.SyncResult{},err
   }
   ref,err:=mailReference(detail)
   if err!=nil{return integrationresponse.SyncResult{},err}
   references=append(references,ref)
  }
  if err:=s.syncRepository.SaveMessages(ctx,id,references);err!=nil{return integrationresponse.SyncResult{},err}
  deletedIDs:=make([]string,0,len(deletes))
  for mid:=range deletes{deletedIDs=append(deletedIDs,mid)}
  if err:=s.syncRepository.MarkDeleted(ctx,id,deletedIDs);err!=nil{return integrationresponse.SyncResult{},err}
  result.Processed+=len(references)
  result.Deleted+=len(deletedIDs)
  page=changes.NextPageToken
  if page==""{
   if changes.HistoryID==""{return integrationresponse.SyncResult{},errors.New("gmail history response missing final cursor")}
   if err:=s.syncRepository.SaveCursor(ctx,id,changes.HistoryID);err!=nil{return integrationresponse.SyncResult{},err}
   result.Cursor=changes.HistoryID
   return result,nil
  }
 }
 return integrationresponse.SyncResult{},ErrSyncTooLarge
}
func mailReference(m gmail.Message)(entity.MailReference,error){
 if strings.TrimSpace(m.ID)==""||strings.TrimSpace(m.ThreadID)==""{return entity.MailReference{},fmt.Errorf("incomplete Gmail message metadata")}
 labels,_:=json.Marshal(m.LabelIDs)
 return entity.MailReference{
  ProviderMessageID:m.ID,ThreadID:m.ThreadID,Direction:gmail.MailDirection(m.LabelIDs),
  HistoryID:m.HistoryID,Labels:string(labels),Deleted:false,
 },nil
}
