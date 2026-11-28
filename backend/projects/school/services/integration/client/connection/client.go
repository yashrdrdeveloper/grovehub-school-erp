package connection

import (
 "bytes"
 "context"
 "encoding/json"
 "errors"
 "net/http"
 "net/url"
 "strings"

 sdk "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/internal/generated/integration"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/internal/platformhttp"
)

var ErrInvalidRequest = errors.New("invalid integration request")
type Client struct { api *sdk.APIClient }

func New(origin *url.URL, client *http.Client) *Client {
 cfg:=sdk.NewConfiguration()
 cfg.Servers=sdk.ServerConfigurations{{URL:strings.TrimRight(origin.String(),"/")}}
 cfg.HTTPClient=platformhttp.BoundedClient(client)
 return &Client{api:sdk.NewAPIClient(cfg)}
}

func validID(id string) bool {
 if len(id)!=36{return false}
 for i,c:=range id {
  if i==8||i==13||i==18||i==23 {if c!='-'{return false};continue}
  if !(c>='0'&&c<='9'||c>='a'&&c<='f'||c>='A'&&c<='F'){return false}
 }
 return true
}

func validAction(action string) bool {
 if len(action)==0||len(action)>100{return false}
 for _,c:=range action {if !(c>='a'&&c<='z'||c>='A'&&c<='Z'||c>='0'&&c<='9'||c=='.'||c=='_'||c=='-'){return false}}
 return true
}

func (c *Client) Business(ctx context.Context, scope,id,provider,action string,input json.RawMessage,proof string)(platformhttp.Result,error){
 ctx=platformhttp.WithIdentityProof(ctx,proof)
 switch scope {
 case "integration.create":
  body,ok:=platformhttp.DecodeBody[sdk.CreateConnectionBody](input)
  if !ok{return platformhttp.Result{},ErrInvalidRequest}
  ctx=platformhttp.WithBody(ctx,input)
  _,resp,err:=c.api.DefaultAPI.CreateConnection(ctx).CreateConnectionBody(body).Execute()
  return platformhttp.SDKResult(resp,err)
 case "integration.get":
  if !validID(id){return platformhttp.Result{},ErrInvalidRequest}
  if len(bytes.TrimSpace(input))!=0{return platformhttp.Result{},ErrInvalidRequest}
  _,resp,err:=c.api.DefaultAPI.GetConnection(ctx,id).Execute()
  return platformhttp.SDKResult(resp,err)
 case "integration.list":
  if len(bytes.TrimSpace(input))!=0{return platformhttp.Result{},ErrInvalidRequest}
  _,resp,err:=c.api.DefaultAPI.ListConnections(ctx).Execute()
  return platformhttp.SDKResult(resp,err)
 case "integration.disconnect":
  if !validID(id){return platformhttp.Result{},ErrInvalidRequest}
  if len(bytes.TrimSpace(input))!=0{return platformhttp.Result{},ErrInvalidRequest}
  _,resp,err:=c.api.DefaultAPI.DisconnectConnection(ctx,id).Execute()
  return platformhttp.SDKResult(resp,err)
 case "integration.oauth.start":
  if !validID(id){return platformhttp.Result{},ErrInvalidRequest}
  if len(bytes.TrimSpace(input))!=0{return platformhttp.Result{},ErrInvalidRequest}
  _,resp,err:=c.api.DefaultAPI.BeginOAuth(ctx,id).Execute()
  return platformhttp.SDKResult(resp,err)
 case "integration.test":
  if !validID(id){return platformhttp.Result{},ErrInvalidRequest}
  body,ok:=platformhttp.DecodeBody[sdk.TestConnectionBody](input)
  if !ok{return platformhttp.Result{},ErrInvalidRequest}
  ctx=platformhttp.WithBody(ctx,input)
  _,resp,err:=c.api.DefaultAPI.TestConnection(ctx,id).TestConnectionBody(body).Execute()
  return platformhttp.SDKResult(resp,err)
 case "integration.tests":
  if !validID(id){return platformhttp.Result{},ErrInvalidRequest}
  if len(bytes.TrimSpace(input))!=0{return platformhttp.Result{},ErrInvalidRequest}
  _,resp,err:=c.api.DefaultAPI.ListConnectionTests(ctx,id).Execute()
  return platformhttp.SDKResult(resp,err)
 case "integration.watch":
  if !validID(id){return platformhttp.Result{},ErrInvalidRequest}
  if len(bytes.TrimSpace(input))!=0{return platformhttp.Result{},ErrInvalidRequest}
  _,resp,err:=c.api.DefaultAPI.WatchConnection(ctx,id).Execute()
  return platformhttp.SDKResult(resp,err)
 case "integration.systemmail":
  body,ok:=platformhttp.DecodeBody[sdk.SystemMailSender](input)
  if !ok{return platformhttp.Result{},ErrInvalidRequest}
  ctx=platformhttp.WithBody(ctx,input)
  _,resp,err:=c.api.DefaultAPI.ConfigureSystemMailSender(ctx).SystemMailSender(body).Execute()
  return platformhttp.SDKResult(resp,err)
 case "integration.action":
  if !validID(id){return platformhttp.Result{},ErrInvalidRequest}
  if !validAction(action){return platformhttp.Result{},ErrInvalidRequest}
  body,ok:=platformhttp.DecodeBody[sdk.ProviderActionBody](input)
  if !ok{return platformhttp.Result{},ErrInvalidRequest}
  ctx=platformhttp.WithBody(ctx,input)
  _,resp,err:=c.api.DefaultAPI.ExecuteProviderAction(ctx,id,action).ProviderActionBody(body).Execute()
  return platformhttp.SDKResult(resp,err)
 case "integration.providers":
  if len(bytes.TrimSpace(input))!=0{return platformhttp.Result{},ErrInvalidRequest}
  _,resp,err:=c.api.DefaultAPI.ListProviders(ctx).Execute()
  return platformhttp.SDKResult(resp,err)
 case "integration.provider":
  if provider!="gmail"&&provider!="outlook"&&provider!="GMAIL"&&provider!="OUTLOOK"{return platformhttp.Result{},ErrInvalidRequest}
  if len(bytes.TrimSpace(input))!=0{return platformhttp.Result{},ErrInvalidRequest}
  _,resp,err:=c.api.DefaultAPI.GetProvider(ctx,provider).Execute()
  return platformhttp.SDKResult(resp,err)
 default:
  return platformhttp.Result{},ErrInvalidRequest
 }
}

// OAuth uses only the signed workload and provider state, never a browser proof.
func(c *Client) FinishOAuth(ctx context.Context,provider,code,state string)(platformhttp.Result,error){
 if (provider!="gmail"&&provider!="outlook")||len(code)==0||len(code)>4096||len(state)!=43{return platformhttp.Result{},ErrInvalidRequest}
 for _,v:=range state{if !(v>='a'&&v<='z'||v>='A'&&v<='Z'||v>='0'&&v<='9'||v=='_'||v=='-'){return platformhttp.Result{},ErrInvalidRequest}}
 ctx=platformhttp.WithIdentityProof(ctx,"")
 _,resp,err:=c.api.DefaultAPI.FinishOAuth(ctx,provider).Code(code).State(state).Execute()
 return platformhttp.SDKResult(resp,err)
}
