package impl

import (
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
 "grovehub/services/integration/providers"
)

const graphBase="https://graph.microsoft.com/v1.0"
type Client struct {
 config integrationconfig.Outlook
 http *http.Client
}
func NewClient(config integrationconfig.Outlook)*Client{
 return &Client{config:config,http:&http.Client{Timeout:20*time.Second}}
}
func(c *Client)endpoint(path string)string{
 return "https://login.microsoftonline.com/"+c.config.Tenant+"/oauth2/v2.0/"+path
}
func(c *Client)AuthorizationURL(state string)string{
 v:=url.Values{
  "client_id":{c.config.ClientID},
  "redirect_uri":{c.config.RedirectURI},
  "response_type":{"code"},
  "response_mode":{"query"},
  "scope":{"offline_access User.Read Mail.Read Mail.Send"},
  "state":{state},
 }
 return c.endpoint("authorize")+"?"+v.Encode()
}
func(c *Client)Exchange(ctx context.Context,code string)(providers.OAuthToken,error){
 return c.tokenRequest(ctx,url.Values{
  "client_id":{c.config.ClientID},
  "client_secret":{c.config.ClientSecret},
  "code":{code},
  "grant_type":{"authorization_code"},
  "redirect_uri":{c.config.RedirectURI},
  "scope":{"offline_access User.Read Mail.Read Mail.Send"},
 })
}
func(c *Client)Refresh(ctx context.Context,t providers.OAuthToken)(providers.OAuthToken,error){
 next,err:=c.tokenRequest(ctx,url.Values{
  "client_id":{c.config.ClientID},
  "client_secret":{c.config.ClientSecret},
  "refresh_token":{t.RefreshToken},
  "grant_type":{"refresh_token"},
  "scope":{"offline_access User.Read Mail.Read Mail.Send"},
 })
 if err!=nil{return next,err}
 if next.RefreshToken==""{next.RefreshToken=t.RefreshToken}
 return next,nil
}
func(c *Client)tokenRequest(ctx context.Context,v url.Values)(providers.OAuthToken,error){
 var token providers.OAuthToken
 req,err:=http.NewRequestWithContext(ctx,http.MethodPost,c.endpoint("token"),strings.NewReader(v.Encode()))
 if err!=nil{return token,err}
 req.Header.Set("Content-Type","application/x-www-form-urlencoded")
 res,err:=c.http.Do(req);if err!=nil{return token,err}
 defer res.Body.Close()
 if res.StatusCode!=http.StatusOK{return token,fmt.Errorf("outlook OAuth token error (HTTP %d)",res.StatusCode)}
 if err:=json.NewDecoder(io.LimitReader(res.Body,32768)).Decode(&token);err!=nil{return token,err}
 if token.AccessToken==""{return token,errors.New("Outlook token response missing access token")}
 token.Expiry=time.Now().Add(time.Duration(token.ExpiresIn)*time.Second)
 return token,nil
}
func(c *Client)graph(ctx context.Context,accessToken,method,path string,body io.Reader,contentType string,dst any)error{
 if !strings.HasPrefix(path,"/"){return errors.New("invalid Graph path")}
 req,err:=http.NewRequestWithContext(ctx,method,graphBase+path,body)
 if err!=nil{return err}
 req.Header.Set("Authorization","Bearer "+accessToken)
 if contentType!=""{req.Header.Set("Content-Type",contentType)}
 res,err:=c.http.Do(req);if err!=nil{return err}
 defer res.Body.Close()
 if res.StatusCode<200||res.StatusCode>=300{return providers.NormalizeStatus(res.StatusCode)}
 if dst==nil{return nil}
 return json.NewDecoder(io.LimitReader(res.Body,12<<20)).Decode(dst)
}
type StatusError struct{StatusCode int}
func(e *StatusError)Error()string{return fmt.Sprintf("Outlook Graph request failed (HTTP %d)",e.StatusCode)}
func(c *Client)Mailbox(ctx context.Context,token providers.OAuthToken)(string,error){
 var me struct{
  Mail string `json:"mail"`
  UserPrincipalName string `json:"userPrincipalName"`
 }
 err:=c.graph(ctx,token.AccessToken,http.MethodGet,"/me?$select=mail,userPrincipalName",nil,"",&me)
 if err!=nil{return "",err}
 email:=me.Mail
 if email==""{email=me.UserPrincipalName}
 if email==""{return "",errors.New("Outlook profile returned no account identity")}
 return strings.ToLower(email),nil
}
