// Package authentication adapts the Platform Identity SDK to School ingress.
package authentication

import (
 "context"
 "fmt"
 "net/http"
 "net/url"
 "strings"

 platformhttp "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/internal/platformhttp"
 sdk "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/internal/generated/identity"
)

type Login struct { Identifier string `json:"identifier"`; Password string `json:"password"` }
type ForgotPassword struct { Email string `json:"email"` }
type Result = platformhttp.Result
type Client struct { api *sdk.APIClient }

func New(origin *url.URL, client *http.Client) *Client {
 cfg := sdk.NewConfiguration()
 cfg.Servers = sdk.ServerConfigurations{{URL: strings.TrimRight(origin.String(), "/")}}
 cfg.HTTPClient = platformhttp.BoundedClient(client)
 return &Client{api:sdk.NewAPIClient(cfg)}
}

var sdkResult = platformhttp.SDKResult

func (c *Client) Login(ctx context.Context, input Login)(Result,error) {
 _,resp,err := c.api.DefaultAPI.Login(WithIdentityProof(ctx,"")).LoginBody(sdk.LoginBody{Identifier:input.Identifier,Password:input.Password}).Execute()
 return sdkResult(resp,err)
}
func (c *Client) ForgotPassword(ctx context.Context,input ForgotPassword)(Result,error) {
 _,resp,err := c.api.DefaultAPI.ForgotPassword(WithIdentityProof(ctx,"")).ForgotBody(sdk.ForgotBody{Email:input.Email}).Execute()
 return sdkResult(resp,err)
}
// Compatibility names for School callers; transport ownership is service-neutral.
var WithIdentityProof = platformhttp.WithIdentityProof
type ProofTransport = platformhttp.ProofTransport
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
