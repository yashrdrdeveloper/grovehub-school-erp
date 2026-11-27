// Package authentication adapts the Platform Identity SDK to School ingress.
// The transitional transport is retained only for business wrappers and Integration OAuth.
package authentication

import (
 "bytes"
 "context"
 "encoding/json"
 "fmt"
 "io"
 "net/http"
 "net/url"
 "strings"

 sdk "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/internal/generated/identity"
)

type Login struct { Identifier string `json:"identifier"`; Password string `json:"password"` }
type ForgotPassword struct { Email string `json:"email"` }
type Result struct { Status int; Body []byte; Cookies []*http.Cookie }
type Client struct { origin *url.URL; http *http.Client; api *sdk.APIClient }

func New(origin *url.URL, client *http.Client) *Client {
 cfg := sdk.NewConfiguration()
 cfg.Servers = sdk.ServerConfigurations{{URL: strings.TrimRight(origin.String(), "/")}}
 bounded := *client
 bounded.Transport = identityResponseTransport{base: client.Transport}
 cfg.HTTPClient = &bounded
 return &Client{origin:origin,http:client,api:sdk.NewAPIClient(cfg)}
}

// Bound and validate responses before the generated SDK consumes them.
// Preserve the original response body and Set-Cookie headers for School ingress.
type identityResponseTransport struct { base http.RoundTripper }
func (t identityResponseTransport) RoundTrip(req *http.Request) (*http.Response,error) {
 base := t.base
 if base == nil { base = http.DefaultTransport }
 resp,err := base.RoundTrip(req)
 if err != nil { return nil,err }
 if resp == nil || resp.Body == nil { return nil,fmt.Errorf("empty platform response") }
 defer resp.Body.Close()
 if resp.StatusCode < 200 || resp.StatusCode >= 600 || (resp.StatusCode >= 300 && resp.StatusCode < 400) { return nil,fmt.Errorf("unexpected platform response") }
 body,err := io.ReadAll(io.LimitReader(resp.Body,(1<<20)+1))
 if err != nil || len(body)>1<<20 || !json.Valid(body) { return nil,fmt.Errorf("invalid platform response") }
 resp.Body = io.NopCloser(bytes.NewReader(body))
 return resp,nil
}
func sdkResult(resp *http.Response, callErr error) (Result,error) {
 if resp == nil { return Result{},fmt.Errorf("platform unavailable") }
 defer resp.Body.Close()
 if resp.StatusCode < 300 && callErr != nil { return Result{},fmt.Errorf("invalid platform success response") }
 body,err := io.ReadAll(io.LimitReader(resp.Body,(1<<20)+1))
 if err != nil || len(body)>1<<20 || !json.Valid(body) { return Result{},fmt.Errorf("invalid platform response") }
 return Result{Status:resp.StatusCode,Body:body,Cookies:resp.Cookies()},nil
}
func (c *Client) Login(ctx context.Context, input Login)(Result,error) {
 _,resp,err := c.api.DefaultAPI.Login(WithIdentityProof(ctx,"")).LoginBody(sdk.LoginBody{Identifier:input.Identifier,Password:input.Password}).Execute()
 return sdkResult(resp,err)
}
func (c *Client) ForgotPassword(ctx context.Context,input ForgotPassword)(Result,error) {
 _,resp,err := c.api.DefaultAPI.ForgotPassword(WithIdentityProof(ctx,"")).ForgotBody(sdk.ForgotBody{Email:input.Email}).Execute()
 return sdkResult(resp,err)
}
func (c *Client) call(ctx context.Context,method,path string,input any,proof string)(Result,error) {
 return c.callQuery(ctx,method,path,input,proof,nil)
}
func (c *Client) callQuery(ctx context.Context,method,path string,input any,proof string,query url.Values)(Result,error) {
 var payload []byte
 if input!=nil {var err error;payload,err=json.Marshal(input);if err!=nil{return Result{},err}}
 ctx=WithIdentityProof(ctx,proof)
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

// WithIdentityProof passes the route-selected credential to the signed transport.
func WithIdentityProof(ctx context.Context, proof string) context.Context {
 return context.WithValue(ctx,proofKey{},proof)
}
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
func (c *Client) Refresh(ctx context.Context, proof string) (Result,error) {
 _,resp,err := c.api.DefaultAPI.Refresh(WithIdentityProof(ctx,proof)).Execute()
 return sdkResult(resp,err)
}
func (c *Client) RecoveryValidate(ctx context.Context, proof string) (Result,error) {
 _,resp,err := c.api.DefaultAPI.ValidateRecovery(WithIdentityProof(ctx,proof)).Execute()
 return sdkResult(resp,err)
}
func (c *Client) InvitationValidate(ctx context.Context, proof string) (Result,error) {
 _,resp,err := c.api.DefaultAPI.ValidateInvitation(WithIdentityProof(ctx,proof)).Execute()
 return sdkResult(resp,err)
}
func (c *Client) RecoveryExchange(ctx context.Context, proof string) (Result,error) {
 _,resp,err := c.api.DefaultAPI.ExchangeRecovery(WithIdentityProof(ctx,proof)).Execute()
 return sdkResult(resp,err)
}
func (c *Client) InvitationExchange(ctx context.Context, proof string) (Result,error) {
 _,resp,err := c.api.DefaultAPI.ExchangeInvitation(WithIdentityProof(ctx,proof)).Execute()
 return sdkResult(resp,err)
}
func (c *Client) ResetStatus(ctx context.Context, proof string) (Result,error) {
 _,resp,err := c.api.DefaultAPI.ResetGrantStatus(WithIdentityProof(ctx,proof)).Execute()
 return sdkResult(resp,err)
}
func (c *Client) ResetPassword(ctx context.Context, input ResetPassword, proof string) (Result,error) {
 _,resp,err := c.api.DefaultAPI.ResetPassword(WithIdentityProof(ctx,proof)).ResetBody(sdk.ResetBody{NewPassword:input.NewPassword}).Execute()
 return sdkResult(resp,err)
}
func (c *Client) SetupStatus(ctx context.Context, proof string) (Result,error) {
 _,resp,err := c.api.DefaultAPI.SetupStatus(WithIdentityProof(ctx,proof)).Execute()
 return sdkResult(resp,err)
}
func (c *Client) Setup(ctx context.Context, input Setup, proof string) (Result,error) {
 _,resp,err := c.api.DefaultAPI.SetupAccount(WithIdentityProof(ctx,proof)).SetupBody(sdk.SetupBody{Password:input.Password}).Execute()
 return sdkResult(resp,err)
}

type ChangePassword struct { NewPassword string `json:"newPassword"` }
type Reauthenticate struct { Password string `json:"password"` }
func (c *Client) Logout(ctx context.Context, proof string) (Result,error) {
 _,resp,err := c.api.DefaultAPI.Logout(WithIdentityProof(ctx,proof)).Execute()
 return sdkResult(resp,err)
}
func (c *Client) LogoutAll(ctx context.Context, proof string) (Result,error) {
 _,resp,err := c.api.DefaultAPI.LogoutAll(WithIdentityProof(ctx,proof)).Execute()
 return sdkResult(resp,err)
}
func (c *Client) ChangePassword(ctx context.Context, input ChangePassword, proof string) (Result,error) {
 _,resp,err := c.api.DefaultAPI.ChangePassword(WithIdentityProof(ctx,proof)).ChangePasswordBody(sdk.ChangePasswordBody{NewPassword:input.NewPassword}).Execute()
 return sdkResult(resp,err)
}
func (c *Client) Reauthenticate(ctx context.Context, input Reauthenticate, proof string) (Result,error) {
 _,resp,err := c.api.DefaultAPI.Reauthenticate(WithIdentityProof(ctx,proof)).ReauthenticationBody(sdk.ReauthenticationBody{Password:input.Password}).Execute()
 return sdkResult(resp,err)
}
func (c *Client) Sessions(ctx context.Context, proof string) (Result,error) {
 _,resp,err := c.api.DefaultAPI.ListSessions(WithIdentityProof(ctx,proof)).Execute()
 return sdkResult(resp,err)
}
func (c *Client) RevokeSession(ctx context.Context, sessionID string, proof string) (Result,error) {
 if !ValidSessionID(sessionID) { return Result{},fmt.Errorf("invalid session ID") }
 _,resp,err := c.api.DefaultAPI.RevokeSession(WithIdentityProof(ctx,proof),sessionID).Execute()
 return sdkResult(resp,err)
}
func (c *Client) Me(ctx context.Context, proof string) (Result,error) {
 _,resp,err := c.api.DefaultAPI.Me(WithIdentityProof(ctx,proof)).Execute()
 return sdkResult(resp,err)
}
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
