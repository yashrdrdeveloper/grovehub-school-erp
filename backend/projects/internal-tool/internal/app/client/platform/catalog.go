// Package platform binds exact, approved operations to a provider-owned SDK.
package platform

import (
 "context"
 "crypto/ed25519"
 "encoding/base64"
 "encoding/json"
 "errors"
 "io"
 "net/http"
 "net/url"
 "strings"

 "github.com/google/uuid"
 core "grovehub/libs/core/gateway"
 sdk "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/generated/integration"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/services/workflow"
)

const catalogPath="/integrations/restricted/providers"
const catalogScope="internal.integration.catalog.read"

type Config struct{
 EnvironmentCode,TemplateID,Origin,KeyID,PrivateKey string
}

func Bindings(c Config)(map[workflow.BindingKey]workflow.Adapter,error){
 if c==(Config{}){return nil,nil}
 id,err:=uuid.Parse(c.TemplateID)
 if err!=nil||id==uuid.Nil||id.String()!=c.TemplateID||c.EnvironmentCode==""||c.KeyID==""{return nil,errors.New("incomplete internal tool catalog binding")}
 origin,err:=url.Parse(c.Origin)
 if err!=nil||origin.Host==""||origin.User!=nil||origin.RawQuery!=""||origin.Fragment!=""||(origin.Path!=""&&origin.Path!="/")||(origin.Scheme!="https"&&!(origin.Scheme=="http"&&(origin.Hostname()=="127.0.0.1"||origin.Hostname()=="::1"))){return nil,errors.New("invalid internal tool Platform origin")}
 key,err:=base64.RawURLEncoding.DecodeString(c.PrivateKey)
 if err!=nil||len(key)!=ed25519.PrivateKeySize{return nil,errors.New("invalid internal tool workload key")}
 transport:=&core.WorkloadTransport{Origin:origin,KeyID:c.KeyID,Issuer:"internal-tool",Subject:"internal-tool",Audience:"platform",PrivateKey:ed25519.PrivateKey(key),Policy:catalogPolicy}
 client,err:=core.NewWorkloadHTTPClient(transport);if err!=nil{return nil,err}
 client.Transport=boundedTransport{base:client.Transport}
 cfg:=sdk.NewConfiguration();cfg.Servers=sdk.ServerConfigurations{{URL:strings.TrimRight(origin.String(),"/")}};cfg.HTTPClient=client
 api:=sdk.NewAPIClient(cfg)
 return map[workflow.BindingKey]workflow.Adapter{{EnvironmentCode:c.EnvironmentCode,TemplateID:c.TemplateID}:func(ctx context.Context,call workflow.Invocation)(int,error){
  if !validInvocation(call){return 0,core.ErrWorkloadDenied}
  _,resp,err:=api.DefaultAPI.ListProviders(ctx).Execute()
  if resp==nil{return 0,errors.New("catalog unavailable")}
  if resp.Body!=nil{defer resp.Body.Close()}
  if resp.StatusCode>=200&&resp.StatusCode<300&&err!=nil{return 0,errors.New("invalid catalog response")}
  return resp.StatusCode,nil
 }},nil
}
func validInvocation(call workflow.Invocation)bool{
 if call.Method!="GET"||call.Path!=catalogPath||call.SchoolID!=nil||len(call.Headers)!=0{return false}
 var in map[string]json.RawMessage
 return json.Unmarshal(call.Input,&in)==nil&&in!=nil&&len(in)==0
}
func catalogPolicy(r *http.Request)(core.WorkloadClaims,error){
 if r==nil||r.URL==nil||r.Method!="GET"||r.URL.Path!=catalogPath||r.URL.RawQuery!=""||r.URL.RawPath!=""||r.Body!=nil&&r.Body!=http.NoBody{return core.WorkloadClaims{},core.ErrWorkloadDenied}
 return core.WorkloadClaims{Scopes:[]string{catalogScope}},nil
}
type boundedTransport struct{base http.RoundTripper}
func(t boundedTransport) RoundTrip(r *http.Request)(*http.Response,error){
 resp,err:=t.base.RoundTrip(r)
 if err!=nil{return nil,err}
 if resp==nil||resp.Body==nil{return nil,errors.New("missing catalog response")}
 defer resp.Body.Close()
 raw,err:=io.ReadAll(io.LimitReader(resp.Body,262145))
 if err!=nil||len(raw)>262144||!json.Valid(raw)||resp.StatusCode<200||resp.StatusCode>=600||(resp.StatusCode>=300&&resp.StatusCode<400){return nil,errors.New("invalid catalog response")}
 resp.Body=io.NopCloser(strings.NewReader(string(raw)))
 return resp,nil
}
