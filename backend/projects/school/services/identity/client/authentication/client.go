// Package authentication owns the transitional typed Platform auth client.
// Replace its wire request types with generated types when the restricted
// OpenAPI specification is verified.
package authentication

import (
 "bytes"
 "context"
 "encoding/json"
 "fmt"
 "io"
 "net/http"
 "net/url"
)

type Login struct { Identifier string `json:"identifier"`; Password string `json:"password"` }
type ForgotPassword struct { Email string `json:"email"` }
type Result struct { Status int; Body []byte; Cookies []*http.Cookie }
type Client struct { origin *url.URL; http *http.Client }

func New(origin *url.URL, client *http.Client) *Client { return &Client{origin:origin,http:client} }
func (c *Client) Login(ctx context.Context, input Login)(Result,error) {
 return c.call(ctx,http.MethodPost,"/identity/restricted/auth/login",input,"")
}
func (c *Client) ForgotPassword(ctx context.Context,input ForgotPassword)(Result,error) {
 return c.call(ctx,http.MethodPost,"/identity/restricted/auth/forgot-password",input,"")
}
func (c *Client) call(ctx context.Context,method,path string,input any,proof string)(Result,error) {
 return c.callQuery(ctx,method,path,input,proof,nil)
}
func (c *Client) callQuery(ctx context.Context,method,path string,input any,proof string,query url.Values)(Result,error) {
 var payload []byte
 if input!=nil {var err error;payload,err=json.Marshal(input);if err!=nil{return Result{},err}}
 ctx=context.WithValue(ctx,proofKey{},proof)
 target:=c.origin.ResolveReference(&url.URL{Path:path,RawQuery:query.Encode()})
 req,err:=http.NewRequestWithContext(ctx,method,target.String(),bytes.NewReader(payload))
 if err!=nil {return Result{},err}
 req.Header.Set("Content-Type","application/json")
 req.Header.Set("Accept","application/json")
 resp,err:=c.http.Do(req)
 if err!=nil {return Result{},err}
 defer resp.Body.Close()
 if resp.StatusCode<200 || resp.StatusCode>=600 || (resp.StatusCode>=300 && resp.StatusCode<400) {return Result{},fmt.Errorf("unexpected platform response")}
 body,err:=io.ReadAll(io.LimitReader(resp.Body,1<<20+1))
 if err!=nil || len(body)>1<<20 || !json.Valid(body) {return Result{},fmt.Errorf("invalid platform response")}
 return Result{Status:resp.StatusCode,Body:body,Cookies:resp.Cookies()},nil
}

type proofKey struct{}
// ProofTransport attaches a single route-selected credential after workload
// signing. It never copies incoming browser headers or cookies.
type ProofTransport struct { Base http.RoundTripper }
func (t ProofTransport) RoundTrip(r *http.Request)(*http.Response,error) {
 clone:=r.Clone(r.Context())
 if proof,ok:=r.Context().Value(proofKey{}).(string);ok && proof!="" {clone.Header.Set("X-GH-Identity-Proof",proof)}
 base:=t.Base
 if base==nil {base=http.DefaultTransport}
 return base.RoundTrip(clone)
}
type ResetPassword struct { NewPassword string `json:"newPassword"` }
type Setup struct { Password string `json:"password"` }
func(c *Client) Refresh(ctx context.Context,proof string)(Result,error){return c.call(ctx,"POST","/identity/restricted/auth/refresh",nil,proof)}
func(c *Client) RecoveryValidate(ctx context.Context,proof string)(Result,error){return c.call(ctx,"GET","/identity/restricted/auth/recovery/validate",nil,proof)}
func(c *Client) InvitationValidate(ctx context.Context,proof string)(Result,error){return c.call(ctx,"GET","/identity/restricted/auth/invitations/validate",nil,proof)}
func(c *Client) RecoveryExchange(ctx context.Context,proof string)(Result,error){return c.call(ctx,"POST","/identity/restricted/auth/recovery/exchange",nil,proof)}
func(c *Client) InvitationExchange(ctx context.Context,proof string)(Result,error){return c.call(ctx,"POST","/identity/restricted/auth/invitations/exchange",nil,proof)}
func(c *Client) ResetStatus(ctx context.Context,proof string)(Result,error){return c.call(ctx,"GET","/identity/restricted/auth/reset-password/status",nil,proof)}
func(c *Client) ResetPassword(ctx context.Context,input ResetPassword,proof string)(Result,error){return c.call(ctx,"POST","/identity/restricted/auth/reset-password",input,proof)}
func(c *Client) SetupStatus(ctx context.Context,proof string)(Result,error){return c.call(ctx,"GET","/identity/restricted/accounts/setup/status",nil,proof)}
func(c *Client) Setup(ctx context.Context,input Setup,proof string)(Result,error){return c.call(ctx,"POST","/identity/restricted/accounts/setup",input,proof)}

type ChangePassword struct { NewPassword string `json:"newPassword"` }
type Reauthenticate struct { Password string `json:"password"` }
func(c *Client) Logout(ctx context.Context,proof string)(Result,error){return c.call(ctx,"POST","/identity/restricted/auth/logout",nil,proof)}
func(c *Client) LogoutAll(ctx context.Context,proof string)(Result,error){return c.call(ctx,"POST","/identity/restricted/auth/logout-all",nil,proof)}
func(c *Client) ChangePassword(ctx context.Context,input ChangePassword,proof string)(Result,error){return c.call(ctx,"POST","/identity/restricted/auth/change-password",input,proof)}
func(c *Client) Reauthenticate(ctx context.Context,input Reauthenticate,proof string)(Result,error){return c.call(ctx,"POST","/identity/restricted/auth/reauthenticate",input,proof)}
func(c *Client) Sessions(ctx context.Context,proof string)(Result,error){return c.call(ctx,"GET","/identity/restricted/auth/sessions",nil,proof)}
func(c *Client) RevokeSession(ctx context.Context,sessionID string,proof string)(Result,error){if !ValidSessionID(sessionID){return Result{},fmt.Errorf("invalid session ID")};return c.call(ctx,"DELETE","/identity/restricted/auth/sessions"+"/"+sessionID,nil,proof)}
func(c *Client) Me(ctx context.Context,proof string)(Result,error){return c.call(ctx,"GET","/identity/restricted/auth/me",nil,proof)}
func ValidSessionID(s string)bool {
 if len(s)!=36{return false}
 for i,c:=range s {
  if i==8||i==13||i==18||i==23 {if c!='-'{return false};continue}
  if !(c>='0'&&c<='9'||c>='a'&&c<='f'||c>='A'&&c<='F'){return false}
 }
 return true
}

// CallBusiness is an internal transport adapter; the configured workload policy
// rejects any method and destination outside the explicit operation allowlist.
func(c *Client) CallBusiness(ctx context.Context,method,path string,input json.RawMessage,proof string)(Result,error){
 if len(input)==0{return c.call(ctx,method,path,nil,proof)}
 return c.call(ctx,method,path,input,proof)
}

func(c *Client) FinishOAuth(ctx context.Context,provider,code,state string)(Result,error){
 if provider!="gmail" && provider!="outlook" {return Result{},fmt.Errorf("invalid provider")}
 return c.callQuery(ctx,http.MethodGet,"/integrations/restricted/connections/oauth/"+provider+"/callback",nil,"",url.Values{"code":{code},"state":{state}})
}
