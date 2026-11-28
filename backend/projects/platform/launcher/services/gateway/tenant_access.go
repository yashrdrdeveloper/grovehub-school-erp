package gateway

import (
 "context"
 "crypto/ed25519"
 "encoding/base64"
 "errors"
 "io"
 "net/http"
 "net/url"
 "strings"
 "strconv"
 "time"

 core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
 "grovehub/libs/schoolgateway"
)

type TenantAccessConfig struct{Origin,KeyID,PrivateKey string}
const tenantAccessSubject="platform-tenant-access"
const tenantAccessScope="platform.school.access-check"

func NewTenantAccess(c TenantAccessConfig)(schoolgateway.TenantCheck,error){
 if c==(TenantAccessConfig{}){return nil,nil}
 origin,err:=url.Parse(c.Origin)
 if err!=nil||origin.Hostname()==""||origin.User!=nil||origin.RawPath!=""||origin.RawQuery!=""||origin.ForceQuery||origin.Fragment!=""||(origin.Path!=""&&origin.Path!="/")||(origin.Scheme!="https"&&!(origin.Scheme=="http"&&(origin.Hostname()=="127.0.0.1"||origin.Hostname()=="::1"))){return nil,errors.New("invalid tenant access origin")}
 private,err:=base64.RawURLEncoding.DecodeString(c.PrivateKey)
 if err!=nil||len(private)!=ed25519.PrivateKeySize||c.KeyID==""{return nil,errors.New("invalid tenant access credential")}
 base:=http.DefaultTransport.(*http.Transport).Clone();base.Proxy=nil
 client,err:=core.NewWorkloadHTTPClient(&core.WorkloadTransport{Origin:origin,KeyID:c.KeyID,Issuer:tenantAccessSubject,Subject:tenantAccessSubject,Audience:"school",PrivateKey:ed25519.PrivateKey(private),Base:base,Policy:func(r *http.Request)(core.WorkloadClaims,error){
  id:=strings.TrimSuffix(strings.TrimPrefix(r.URL.Path,"/school/restricted/registry/schools/"),"/access-check")
  if !validSchoolID(id)||r.URL.Path!="/school/restricted/registry/schools/"+id+"/access-check"||r.Method!="GET"||r.URL.RawQuery!=""||r.URL.RawPath!=""||r.URL.ForceQuery||r.ContentLength!=0{return core.WorkloadClaims{},core.ErrWorkloadDenied}
  return core.WorkloadClaims{Scopes:[]string{tenantAccessScope}},nil
 }})
 if err!=nil{return nil,err}
 return func(ctx context.Context,id string)(int64,error){
  if !validSchoolID(id){return 0,schoolgateway.ErrTenantDenied}
  ctx,cancel:=context.WithTimeout(ctx,3*time.Second);defer cancel()
  req,err:=http.NewRequestWithContext(ctx,"GET",strings.TrimRight(origin.String(),"/")+"/school/restricted/registry/schools/"+id+"/access-check",nil)
  if err!=nil{return 0,schoolgateway.ErrTenantDenied}
  resp,err:=client.Do(req);if err!=nil{return 0,schoolgateway.ErrTenantDenied};defer resp.Body.Close()
  body,err:=io.ReadAll(io.LimitReader(resp.Body,1))
  if err!=nil||resp.StatusCode!=204||len(body)!=0{return 0,schoolgateway.ErrTenantDenied}
  values:=resp.Header.Values("X-GH-Session-Generation")
  if len(values)!=1{return 0,schoolgateway.ErrTenantDenied}
  generation,err:=strconv.ParseInt(values[0],10,64)
  if err!=nil||generation<1||strconv.FormatInt(generation,10)!=values[0]{return 0,schoolgateway.ErrTenantDenied}
  return generation,nil
 },nil
}
