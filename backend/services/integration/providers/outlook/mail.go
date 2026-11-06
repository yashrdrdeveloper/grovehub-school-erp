package outlook

import (
 "context"
 "encoding/base64"
 "encoding/json"
 "net/http"
 "net/url"
 "strings"

 "grovehub/services/integration/providers"
)

type Message struct{
 ID string `json:"id"`
 Subject string `json:"subject"`
 ConversationID string `json:"conversationId"`
 InternetMessageID string `json:"internetMessageId"`
 ReceivedDateTime string `json:"receivedDateTime"`
 SentDateTime string `json:"sentDateTime"`
 From json.RawMessage `json:"from,omitempty"`
 ToRecipients json.RawMessage `json:"toRecipients,omitempty"`
 Body json.RawMessage `json:"body,omitempty"`
 IsRead bool `json:"isRead"`
}
type MessagePage struct{
 Messages []Message `json:"value"`
 NextLink string `json:"@odata.nextLink,omitempty"`
}
type SendResult struct{Accepted bool `json:"accepted"`}

func(c *Client)Profile(ctx context.Context,accessToken string)(any,error){
 var result struct{
  Mail string `json:"mail"`
  UserPrincipalName string `json:"userPrincipalName"`
 }
 err:=c.graph(ctx,accessToken,http.MethodGet,"/me?$select=mail,userPrincipalName",nil,"",&result)
 if err!=nil{return nil,err}
 if result.Mail==""{result.Mail=result.UserPrincipalName}
 return map[string]string{"emailAddress":strings.ToLower(result.Mail)},nil
}
func safeNextLink(raw string)(string,error){
 if raw==""{return "/me/messages?$top=50&$select=id,subject,conversationId,internetMessageId,receivedDateTime,sentDateTime,from,toRecipients,isRead",nil}
 if len(raw)>4096{return "",providers.ErrInvalidActionInput}
 parsed,err:=url.Parse(raw);if err!=nil{return "",providers.ErrInvalidActionInput}
 if parsed.Scheme!="https"||parsed.Host!="graph.microsoft.com"||parsed.User!=nil||parsed.Fragment!=""||
  parsed.Path!="/v1.0/me/messages" {
  return "",providers.ErrInvalidActionInput
 }
 return "/me/messages?"+parsed.RawQuery,nil
}
func(c *Client)Messages(ctx context.Context,accessToken,page string)(MessagePage,error){
 var result MessagePage
 path,err:=safeNextLink(page);if err!=nil{return result,err}
 err=c.graph(ctx,accessToken,http.MethodGet,path,nil,"",&result)
 return result,err
}
func(c *Client)Message(ctx context.Context,accessToken,id string)(Message,error){
 var result Message
 if id==""||len(id)>512||strings.ContainsAny(id,"?#/"){return result,providers.ErrInvalidActionInput}
 err:=c.graph(ctx,accessToken,http.MethodGet,"/me/messages/"+url.PathEscape(id),nil,"",&result)
 return result,err
}
func(c *Client)SendRaw(ctx context.Context,accessToken,raw string)(SendResult,error){
 var result SendResult
 if len(raw)==0||len(raw)>8*1024*1024{return result,providers.ErrInvalidActionInput}
 decoded,err:=base64.RawURLEncoding.DecodeString(raw)
 if err!=nil{return result,providers.ErrInvalidActionInput}
 if len(decoded)==0||len(decoded)>6*1024*1024{return result,providers.ErrInvalidActionInput}
 // Microsoft Graph's MIME form expects base64 (not Gmail's base64url)
 // supplied as text/plain; response is HTTP 202 with no message identifier.
 encoded:=base64.StdEncoding.EncodeToString(decoded)
 err=c.graph(ctx,accessToken,http.MethodPost,"/me/sendMail",strings.NewReader(encoded),"text/plain",nil)
 if err!=nil{return result,err}
 result.Accepted=true
 return result,nil
}
