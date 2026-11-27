// Package audit wraps the provider-owned generated Audit API client.
package audit

import (
 "context"
 "crypto/ed25519"
 "encoding/base64"
 "errors"
 "net/http"
 "net/url"
 "strings"
 "time"

 core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
 sdk "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/generated/audit"
)

const path="/audit/restricted/operator-activities"
const scope="internal.audit.operator.write"

type Config struct{ Origin,KeyID,PrivateKey string }
type Activity struct {
 EventID,ActorID,ResourceID,Operation string
 OccurredAt time.Time
}
type Publisher func(context.Context,Activity)(int,error)

func buildClient(c Config,subject string,authorize func(*http.Request)(core.WorkloadClaims,error))(*sdk.APIClient,error){
 if c==(Config{}){return nil,nil}
 origin,err:=url.Parse(c.Origin)
 if err!=nil||origin.Hostname()==""||origin.User!=nil||origin.RawQuery!=""||origin.ForceQuery||origin.Fragment!=""||origin.RawPath!=""||(origin.Path!=""&&origin.Path!="/")||(origin.Scheme!="https"&&!(origin.Scheme=="http"&&(origin.Hostname()=="127.0.0.1"||origin.Hostname()=="::1"))){return nil,errors.New("invalid audit origin")}
 key,err:=base64.RawURLEncoding.DecodeString(c.PrivateKey)
 if err!=nil||len(key)!=ed25519.PrivateKeySize||c.KeyID==""{return nil,errors.New("invalid audit workload key")}
 transport:=&core.WorkloadTransport{Origin:origin,KeyID:c.KeyID,Issuer:subject,Subject:subject,Audience:"platform",PrivateKey:ed25519.PrivateKey(key),Policy:authorize}
 httpClient,err:=core.NewWorkloadHTTPClient(transport);if err!=nil{return nil,err}
 cfg:=sdk.NewConfiguration();cfg.Servers=sdk.ServerConfigurations{{URL:strings.TrimRight(origin.String(),"/")}};cfg.HTTPClient=httpClient
 return sdk.NewAPIClient(cfg),nil
}
func New(c Config)(Publisher,error){
 api,err:=buildClient(c,"internal-tool-audit",policy)
 if err!=nil||api==nil{return nil,err}
 return func(ctx context.Context,a Activity)(int,error){
  body:=sdk.OperatorActivity{EventId:a.EventID,ActorId:a.ActorID,ResourceId:a.ResourceID,Operation:a.Operation,OccurredAt:a.OccurredAt}
  resp,err:=api.DefaultAPI.RecordOperatorActivity(ctx).OperatorActivity(body).Execute()
  if resp==nil{return 0,errors.New("audit unavailable")}
  if resp.Body!=nil{defer resp.Body.Close()}
  if resp.StatusCode==204&&err!=nil{return 0,errors.New("invalid audit acknowledgment")}
  return resp.StatusCode,nil
 },nil
}
func policy(r *http.Request)(core.WorkloadClaims,error){
 if r==nil||r.URL==nil||r.Method!="POST"||r.URL.Path!=path||r.URL.RawPath!=""||r.URL.RawQuery!=""||r.URL.ForceQuery||r.Header.Get("Cookie")!=""||r.Header.Get("X-GH-Identity-Proof")!=""{return core.WorkloadClaims{},core.ErrWorkloadDenied}
 return core.WorkloadClaims{Scopes:[]string{scope}},nil
}
