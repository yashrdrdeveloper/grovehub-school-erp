package platform

import (
 "context"
 "crypto/ed25519"
 "encoding/base64"
 "encoding/json"
 "errors"
 "net/http"
 "net/url"
 "io"
 "regexp"
 "strings"

 "github.com/google/uuid"
 core "grovehub/libs/core/gateway"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/services/workflow"
 sdk "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/generated/audit"
)

const schoolCheckScope="internal.audit.school.activity-check"
const schoolIssuer="internal-tool-school"
type SchoolConfig struct{
 EnvironmentCode,TemplateID,SchoolID,Origin,KeyID,PrivateKey,OperatorIDs string
}
func canonicalID(v string)bool{id,err:=uuid.Parse(v);return err==nil&&id!=uuid.Nil&&id.String()==v}
func SchoolBindings(c SchoolConfig)(map[workflow.SchoolBindingKey]workflow.SchoolBinding,error){
 if c==(SchoolConfig{}){return nil,nil}
 if !canonicalID(c.TemplateID)||!canonicalID(c.SchoolID)||!regexp.MustCompile("^[a-z][a-z0-9_-]{1,63}$").MatchString(c.EnvironmentCode)||c.KeyID==""{return nil,errors.New("invalid school execution binding")}
 operators:=strings.Split(c.OperatorIDs,",")
 if len(operators)>100{return nil,errors.New("too many school execution operators")}
 seen:=map[string]bool{}
 for _,id:=range operators{if !canonicalID(id)||seen[id]{return nil,errors.New("invalid school execution operator allowlist")};seen[id]=true}
 origin,err:=url.Parse(c.Origin)
 if err!=nil||origin.Hostname()==""||origin.User!=nil||origin.RawQuery!=""||origin.ForceQuery||origin.Fragment!=""||origin.RawPath!=""||(origin.Path!=""&&origin.Path!="/")||(origin.Scheme!="https"&&!(origin.Scheme=="http"&&(origin.Hostname()=="127.0.0.1"||origin.Hostname()=="::1"))){return nil,errors.New("invalid school execution origin")}
 key,err:=base64.RawURLEncoding.DecodeString(c.PrivateKey)
 if err!=nil||len(key)!=ed25519.PrivateKeySize{return nil,errors.New("invalid school execution workload key")}
 path:="/audit/restricted/schools/"+c.SchoolID+"/activity-check"
 validate:=schoolCheckInvocation(c.SchoolID,path)
 policy:=func(r *http.Request)(core.WorkloadClaims,error){
  if r==nil||r.URL==nil||r.Method!="GET"||r.URL.Path!=path||r.URL.RawPath!=""||r.URL.RawQuery!=""||r.URL.ForceQuery||r.ContentLength!=0||len(r.Header.Values("Cookie"))!=0||len(r.Header.Values("X-GH-Identity-Proof"))!=0{return core.WorkloadClaims{},core.ErrWorkloadDenied}
  return core.WorkloadClaims{Scopes:[]string{schoolCheckScope}},nil
 }
 transport:=&core.WorkloadTransport{Origin:origin,KeyID:c.KeyID,Issuer:schoolIssuer,Subject:schoolIssuer+":"+c.SchoolID,Audience:"platform",PrivateKey:ed25519.PrivateKey(key),Policy:policy}
 client,err:=core.NewWorkloadHTTPClient(transport);if err!=nil{return nil,err}
 cfg:=sdk.NewConfiguration();cfg.Servers=sdk.ServerConfigurations{{URL:strings.TrimRight(origin.String(),"/")}};client.Transport=schoolResponseTransport{base:client.Transport}
 cfg.HTTPClient=client
 api:=sdk.NewAPIClient(cfg)
 adapter:=func(ctx context.Context,call workflow.Invocation)(int,error){
  if !validate(call){return 0,core.ErrWorkloadDenied}
  resp,err:=api.DefaultAPI.CheckSchoolAuditActivity(ctx,c.SchoolID).Execute()
  if resp==nil{return 0,errors.New("school audit unavailable")}
  if resp.Body!=nil{defer resp.Body.Close()}
  if resp.StatusCode==204&&err!=nil{return 0,errors.New("invalid school audit acknowledgment")}
  if resp.StatusCode!=204&&resp.StatusCode!=403&&resp.StatusCode!=404&&resp.StatusCode!=503{return 0,errors.New("unexpected school audit status")}
  return resp.StatusCode,nil
 }
 return map[workflow.SchoolBindingKey]workflow.SchoolBinding{{EnvironmentCode:c.EnvironmentCode,TemplateID:c.TemplateID,SchoolID:c.SchoolID}:{Adapter:adapter,Validate:validate,Operators:operators}},nil
}

func schoolCheckInvocation(school,path string)func(workflow.Invocation)bool{
 return func(call workflow.Invocation)bool{
  if call.SchoolID==nil||*call.SchoolID!=school||call.Method!="GET"||call.Path!=path||len(call.Headers)!=0{return false}
  var input map[string]json.RawMessage
  if json.Unmarshal(call.Input,&input)!=nil||len(input)!=1{return false}
  var value string
  return json.Unmarshal(input["schoolId"],&value)==nil&&value==school
 }
}

type schoolResponseTransport struct{base http.RoundTripper}
func(t schoolResponseTransport) RoundTrip(r *http.Request)(*http.Response,error){
 resp,err:=t.base.RoundTrip(r);if err!=nil{return nil,err}
 if resp==nil||resp.Body==nil{return nil,errors.New("missing school audit response")}
 defer resp.Body.Close()
 raw,err:=io.ReadAll(io.LimitReader(resp.Body,4097))
 if err!=nil||len(raw)>4096||(resp.StatusCode==204&&len(raw)!=0){return nil,errors.New("invalid school audit response")}
 resp.Body=io.NopCloser(strings.NewReader(string(raw)))
 return resp,nil
}
