package gmail

import (
 "context"
 "errors"
 "fmt"
 "encoding/base64"
 "net/http"
 "net/url"
 "strings"
)

type MailProfile struct {EmailAddress string `json:"emailAddress"`;HistoryID string `json:"historyId"`}
type SentMessage struct {ID string `json:"id"`; ThreadID string `json:"threadId"`}

type Message struct {
 ID string `json:"id"`
 ThreadID string `json:"threadId"`
 HistoryID string `json:"historyId"`
 LabelIDs []string `json:"labelIds"`
}
type MessagePage struct{
 Messages []Message `json:"messages"`
 NextPageToken string `json:"nextPageToken"`
}
type HistoryChange struct{
 MessagesAdded []struct{Message Message `json:"message"`} `json:"messagesAdded"`
 MessagesDeleted []struct{Message Message `json:"message"`} `json:"messagesDeleted"`
 LabelsAdded []struct{Message Message `json:"message"`} `json:"labelsAdded"`
 LabelsRemoved []struct{Message Message `json:"message"`} `json:"labelsRemoved"`
}
type HistoryPage struct{
 History []HistoryChange `json:"history"`
 NextPageToken string `json:"nextPageToken"`
 HistoryID string `json:"historyId"`
}
var ErrExpiredHistory = errors.New("gmail history cursor expired")

func(c *Client)HistoryID(ctx context.Context,t Token)(string,error){
 var p struct{HistoryID string `json:"historyId"`}
 if err:=c.do(ctx,t,http.MethodGet,"/profile",nil,&p);err!=nil{return "",err}
 if p.HistoryID==""{return "",errors.New("gmail profile has no history cursor")}
 return p.HistoryID,nil
}
func(c *Client)Messages(ctx context.Context,t Token,page string)(MessagePage,error){
 var p MessagePage
 query:=url.Values{"maxResults":{"100"}}
 if page!=""{query.Set("pageToken",page)}
 err:=c.do(ctx,t,http.MethodGet,"/messages?"+query.Encode(),nil,&p)
 return p,err
}
func(c *Client)MessageMetadata(ctx context.Context,t Token,id string)(Message,error){
 var m Message
 // Google returns labelIds, threadId, historyId in METADATA format without body.
 err:=c.do(ctx,t,http.MethodGet,"/messages/"+url.PathEscape(id)+"?format=metadata",nil,&m)
 return m,err
}
func(c *Client)History(ctx context.Context,t Token,startID,page string)(HistoryPage,error){
 var p HistoryPage
 if startID==""{return p,errors.New("missing Gmail startHistoryId")}
 query:=url.Values{"startHistoryId":{startID},"maxResults":{"100"}}
 if page!=""{query.Set("pageToken",page)}
 path:="/history?"+query.Encode()
 err:=c.do(ctx,t,http.MethodGet,path,nil,&p)
 if err!=nil {
  var failure *StatusError
  if errors.As(err,&failure)&&failure.StatusCode==http.StatusNotFound{return p,ErrExpiredHistory}
  return p,err
 }
 return p,nil
}
// SendRaw submits an RFC 5322 message encoded as base64url.
// The caller (eventually Communication) owns its message and delivery state.
func(c *Client)SendRaw(ctx context.Context,t Token,raw string)(SentMessage,error){
 var sent SentMessage
 if len(raw)==0||len(raw)>8*1024*1024{return sent,errors.New("invalid raw email size")}
 if _,err:=base64.RawURLEncoding.DecodeString(raw);err!=nil{return sent,errors.New("invalid base64url email")}
 err:=c.do(ctx,t,http.MethodPost,"/messages/send",map[string]string{"raw":raw},&sent)
 return sent,err
}

func MailDirection(labels []string)string{
 for _,label:=range labels{if strings.EqualFold(label,"SENT"){return "OUTBOUND"}}
 return "INBOUND"
}
func validateMessage(m Message)error{
 if m.ID==""||m.ThreadID==""{return fmt.Errorf("incomplete Gmail message metadata")}
 return nil
}
