package gmail

import (
 "bytes"
 "context"
 "encoding/json"
 "errors"
 "fmt"
 "io"
 "net/http"
 "net/url"
 "strings"
 "time"

 integrationconfig "grovehub/services/integration/config"
)

const (
 authURL="https://accounts.google.com/o/oauth2/v2/auth"
 tokenURL="https://oauth2.googleapis.com/token"
 gmailAPI="https://gmail.googleapis.com/gmail/v1/users/me"
)

type Token struct{
 AccessToken string `json:"access_token"`
 RefreshToken string `json:"refresh_token"`
 ExpiresIn int `json:"expires_in"`
 TokenType string `json:"token_type"`
 Expiry time.Time `json:"expiry"`
}
type Client struct{
 config integrationconfig.Gmail
 http *http.Client
}
func NewClient(config integrationconfig.Gmail)*Client{
 return &Client{config:config,http:&http.Client{Timeout:20*time.Second}}
}
func (c *Client)AuthorizationURL(state string)string{
 query:=url.Values{
 "client_id":{c.config.ClientID},
 "redirect_uri":{c.config.RedirectURI},
 "response_type":{"code"},
 "scope":{"https://www.googleapis.com/auth/gmail.send https://www.googleapis.com/auth/gmail.readonly https://www.googleapis.com/auth/userinfo.email"},
 "access_type":{"offline"},
 "prompt":{"consent"},
 "state":{state},
 }
 return authURL+"?"+query.Encode()
}
func (c *Client)Exchange(ctx context.Context,code string)(Token,error){
 return c.tokenRequest(ctx,url.Values{
 "code":{code},"client_id":{c.config.ClientID},
 "client_secret":{c.config.ClientSecret},"redirect_uri":{c.config.RedirectURI},
 "grant_type":{"authorization_code"},
 })
}
func (c *Client)Refresh(ctx context.Context,token Token)(Token,error){
 next,err:=c.tokenRequest(ctx,url.Values{
 "refresh_token":{token.RefreshToken},"client_id":{c.config.ClientID},
 "client_secret":{c.config.ClientSecret},"grant_type":{"refresh_token"},
 })
 if err!=nil{return Token{},err}
 if next.RefreshToken==""{next.RefreshToken=token.RefreshToken}
 return next,nil
}
func (c *Client)tokenRequest(ctx context.Context,v url.Values)(Token,error){
 req,err:=http.NewRequestWithContext(ctx,http.MethodPost,tokenURL,strings.NewReader(v.Encode()))
 if err!=nil{return Token{},err}
 req.Header.Set("Content-Type","application/x-www-form-urlencoded")
 resp,err:=c.http.Do(req)
 if err!=nil{return Token{},err}
 defer resp.Body.Close()
 if resp.StatusCode!=http.StatusOK{
  return Token{},fmt.Errorf("google token exchange failed (HTTP %d)",resp.StatusCode)
 }
 var token Token
 if err:=json.NewDecoder(io.LimitReader(resp.Body,32768)).Decode(&token);err!=nil{return Token{},err}
 if token.AccessToken==""{return Token{},errors.New("google did not provide access token")}
 token.Expiry=time.Now().Add(time.Duration(token.ExpiresIn)*time.Second)
 return token,nil
}
func(c *Client)do(ctx context.Context,token Token,method,path string,body any,out any)error{
 var reader io.Reader
 if body!=nil {
  payload,err:=json.Marshal(body);if err!=nil{return err}
  reader=bytes.NewReader(payload)
 }
 req,err:=http.NewRequestWithContext(ctx,method,gmailAPI+path,reader)
 if err!=nil{return err}
 req.Header.Set("Authorization","Bearer "+token.AccessToken)
 if body!=nil{req.Header.Set("Content-Type","application/json")}
 resp,err:=c.http.Do(req)
 if err!=nil{return err}
 defer resp.Body.Close()
 if resp.StatusCode<200||resp.StatusCode>=300{return fmt.Errorf("gmail request failed (HTTP %d)",resp.StatusCode)}
 if out==nil{return nil}
 return json.NewDecoder(io.LimitReader(resp.Body,1<<20)).Decode(out)
}
func(c *Client)Mailbox(ctx context.Context,token Token)(string,error){
 var response struct{EmailAddress string `json:"emailAddress"`}
 err:=c.do(ctx,token,http.MethodGet,"/profile",nil,&response)
 if err!=nil{return "",err}
 if response.EmailAddress==""{return "",errors.New("gmail returned empty email address")}
 return strings.ToLower(response.EmailAddress),nil
}
func(c *Client)InboundAccess(ctx context.Context,token Token)error{
 var result struct{ResultSizeEstimate int `json:"resultSizeEstimate"`}
 // A mailbox list probe confirms Gmail read permission, not delivery of an email.
 return c.do(ctx,token,http.MethodGet,"/messages?maxResults=1",nil,&result)
}
func(c *Client)Outbound(ctx context.Context,token Token,raw string)error{
 return c.do(ctx,token,http.MethodPost,"/messages/send",map[string]string{"raw":raw},nil)
}
