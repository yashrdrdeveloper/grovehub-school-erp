package gmail

import (
 "context"
 "errors"
 "fmt"
 "net/http"
 "net/url"
 "strings"
)

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
func MailDirection(labels []string)string{
 for _,label:=range labels{if strings.EqualFold(label,"SENT"){return "OUTBOUND"}}
 return "INBOUND"
}
func validateMessage(m Message)error{
 if m.ID==""||m.ThreadID==""{return fmt.Errorf("incomplete Gmail message metadata")}
 return nil
}
