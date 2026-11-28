package platform

import (
 "context"
 corehttp "grovehub/libs/core/http"
 "crypto/ed25519"
 "encoding/base64"
 "encoding/json"
 "errors"
 "net/http"
 "net/url"
 "regexp"
 "strings"

 core "grovehub/libs/core/gateway"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/services/workflow"
 sdk "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/generated/registry"
)
const statusScope="internal.school.registry.status"
const statusIssuer="internal-tool-school-status"
func validRegistryStatus(s string)bool{return s=="ACTIVE"||s=="DISABLED"}
func statusInput(call workflow.Invocation,school,path string)(sdk.SchoolStatusInput,bool){
 var body sdk.SchoolStatusInput
 if call.SchoolID==nil||*call.SchoolID!=school||call.Method!="PATCH"||call.Path!=path||len(call.Headers)!=0{return body,false}
 var raw map[string]json.RawMessage
 if json.Unmarshal(call.Input,&raw)!=nil||len(raw)!=4{return body,false}
 var id string
 if json.Unmarshal(raw["schoolId"],&id)!=nil||id!=school||json.Unmarshal(raw["expectedStatus"],&body.ExpectedStatus)!=nil||json.Unmarshal(raw["status"],&body.Status)!=nil||json.Unmarshal(raw["expectedVersion"],&body.ExpectedVersion)!=nil{return body,false}
 return body,body.ExpectedVersion>0&&validRegistryStatus(body.ExpectedStatus)&&validRegistryStatus(body.Status)&&body.ExpectedStatus!=body.Status
}

func StatusBindings(c SchoolConfig)(map[workflow.SchoolBindingKey]workflow.SchoolBinding,error){
 if !canonicalID(c.TemplateID)||!canonicalID(c.SchoolID)||!regexp.MustCompile("^[a-z][a-z0-9_-]{1,63}$").MatchString(c.EnvironmentCode)||c.KeyID==""{return nil,errors.New("invalid status binding")}
 operators:=strings.Split(c.OperatorIDs,",");if len(operators)>100{return nil,errors.New("too many status operators")}
 for _,id:=range operators{if !canonicalID(id){return nil,errors.New("invalid status operator")}}
 origin,err:=url.Parse(c.Origin)
 if err!=nil||origin.Hostname()==""||origin.User!=nil||origin.RawQuery!=""||origin.ForceQuery||origin.Fragment!=""||origin.RawPath!=""||(origin.Path!=""&&origin.Path!="/")||(origin.Scheme!="https"&&!(origin.Scheme=="http"&&(origin.Hostname()=="127.0.0.1"||origin.Hostname()=="::1"))){return nil,errors.New("invalid status origin")}
 key,err:=base64.RawURLEncoding.DecodeString(c.PrivateKey);if err!=nil||len(key)!=ed25519.PrivateKeySize{return nil,errors.New("invalid status key")}
 path:="/school/restricted/registry/schools/"+c.SchoolID+"/status"
 validate:=func(call workflow.Invocation)bool{_,ok:=statusInput(call,c.SchoolID,path);return ok}
 base:=http.DefaultTransport.(*http.Transport).Clone();base.Proxy=nil
 transport:=&core.WorkloadTransport{Origin:origin,KeyID:c.KeyID,Issuer:statusIssuer,Subject:statusIssuer+":"+c.SchoolID,Audience:"school",PrivateKey:ed25519.PrivateKey(key),Base:base,Policy:func(r *http.Request)(core.WorkloadClaims,error){
  if r==nil||r.URL==nil||r.Method!="PATCH"||r.URL.Path!=path||r.URL.RawPath!=""||r.URL.RawQuery!=""||r.URL.ForceQuery||len(r.Header.Values("Idempotency-Key"))!=1||!canonicalID(r.Header.Get("Idempotency-Key"))||len(r.Header.Values("Cookie"))!=0||len(r.Header.Values("X-GH-Identity-Proof"))!=0{return core.WorkloadClaims{},core.ErrWorkloadDenied}
  return core.WorkloadClaims{Scopes:[]string{statusScope}},nil
 }}
 client,err:=core.NewWorkloadHTTPClient(transport);if err!=nil{return nil,err}
 client.Transport=schoolResponseTransport{base:client.Transport}
 cfg:=sdk.NewConfiguration();cfg.Servers=sdk.ServerConfigurations{{URL:strings.TrimRight(origin.String(),"/")}};cfg.HTTPClient=client
 api:=sdk.NewAPIClient(cfg)
 adapter:=func(ctx context.Context,call workflow.Invocation)(int,error){
  body,ok:=statusInput(call,c.SchoolID,path);if !ok||!canonicalID(call.ExecutionID){return 0,core.ErrWorkloadDenied}
  if id:=corehttp.ID(ctx);id!=""{body.CorrelationId=&id}
  resp,err:=api.DefaultAPI.ChangeSchoolStatus(ctx,c.SchoolID).IdempotencyKey(call.ExecutionID).SchoolStatusInput(body).Execute()
  if resp==nil{return 0,errors.New("status unavailable")}
  if resp.Body!=nil{defer resp.Body.Close()}
  if resp.StatusCode==204&&err!=nil{return 0,errors.New("invalid status acknowledgment")}
  switch resp.StatusCode{case 204,400,403,404,409,503:return resp.StatusCode,nil}
  return 0,errors.New("unexpected status status")
 }
 return map[workflow.SchoolBindingKey]workflow.SchoolBinding{{EnvironmentCode:c.EnvironmentCode,TemplateID:c.TemplateID,SchoolID:c.SchoolID}:{Adapter:adapter,Validate:validate,Operators:operators,Mutation:true,Lifecycle:true}},nil
}
