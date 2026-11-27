package gateway

import (
 "context"
 "crypto/ed25519"
 "crypto/rand"
 "encoding/base64"
 "encoding/json"
 "errors"
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"
 "time"

 "github.com/google/uuid"
 core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
 school "grovehub/libs/schoolgateway"
)
func TestTenantAccessClient(t *testing.T){
 id:=uuid.NewString()
 pub,priv,err:=ed25519.GenerateKey(rand.Reader);if err!=nil{t.Fatal(err)}
 keys,err:=core.NewEd25519Credentials([]core.WorkloadKeyConfig{{KeyID:"access",PublicKey:base64.RawURLEncoding.EncodeToString(pub),Issuer:tenantAccessSubject,Subject:tenantAccessSubject,Scopes:[]string{tenantAccessScope}}});if err!=nil{t.Fatal(err)}
 verifier:=core.ScopedWorkloadVerifier{Credentials:keys,Audience:"school",Routes:map[string]core.RoutePolicy{"GET /school/restricted/registry/schools/{schoolId}/access-check":{Scope:tenantAccessScope}}}
 status:=204
 server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  r.Pattern="GET /school/restricted/registry/schools/{schoolId}/access-check"
  if _,err:=verifier.VerifyWorkload(r.Context(),r);err!=nil||r.URL.Path!="/school/restricted/registry/schools/"+id+"/access-check"{t.Error("invalid signed request",err);w.WriteHeader(403);return}
  w.WriteHeader(status)
 }));defer server.Close()
 check,err:=NewTenantAccess(TenantAccessConfig{Origin:server.URL,KeyID:"access",PrivateKey:base64.RawURLEncoding.EncodeToString(priv)});if err!=nil{t.Fatal(err)}
 for _,code:=range []int{204,403,503,302,200}{
  status=code
  if err=check(context.Background(),id);(err==nil)!=(code==204){t.Fatal("incorrect decision",code,err)}
 }
 if err=check(context.Background(),"invalid");err==nil{t.Fatal("invalid ID accepted")}
 if _,err=NewTenantAccess(TenantAccessConfig{Origin:server.URL});err==nil{t.Fatal("partial configuration accepted")}
}
type tenantIdentity struct{tenant string}
func(i tenantIdentity) VerifyAccess(context.Context,string)(school.Principal,error){return school.Principal{TenantID:i.tenant,AccountID:"account",PersonGID:"person",SessionID:"session"},nil}
func(i tenantIdentity) VerifyDelegatedSession(context.Context,string)(school.Principal,error){return school.Principal{},errors.New("not used")}
func(i tenantIdentity) AuthorizeProtectedRoute(context.Context,string)error{return nil}
func TestActorTenantEnforcement(t *testing.T){
 pub,priv,err:=ed25519.GenerateKey(rand.Reader);if err!=nil{t.Fatal(err)}
 id:=uuid.NewString()
 scopes:=[]string{"identity.auth.me","identity.auth.logout","identity.auth.login"}
 key:=core.WorkloadKeyConfig{KeyID:"school",PublicKey:base64.RawURLEncoding.EncodeToString(pub),Issuer:"school",Subject:"school",Scopes:scopes}
 routes:=map[string]core.RoutePolicy{}
 for _,pattern:=range []string{"GET /identity/restricted/auth/me","POST /identity/restricted/auth/logout","POST /identity/restricted/auth/login"}{routes[pattern]=restrictedPolicies[pattern]}
 raw,_:=json.Marshal(workloadConfig{Keys:[]core.WorkloadKeyConfig{key},Routes:routes})
 active:=false
 check:=func(_ context.Context,tenant string)error{
  if tenant!=id{t.Fatal("untrusted tenant used",tenant)}
  if !active{return school.ErrTenantDenied};return nil
 }
 v,err:=NewPlatformWorkloads(string(raw),tenantIdentity{id},nil,check);if err!=nil{t.Fatal(err)}
 call:=func(pattern string)(context.Context,error){
  policy:=routes[pattern]
  token,err:=core.SignWorkload(key.KeyID,key.Issuer,priv,core.WorkloadClaims{Subject:key.Subject,Audience:"platform",Scopes:[]string{policy.Scope},ExpiresAt:time.Now().Add(time.Minute)});if err!=nil{t.Fatal(err)}
  pieces:=strings.SplitN(pattern," ",2)
  r:=httptest.NewRequest(pieces[0],pieces[1],nil);r.Pattern=pattern;r.Header.Set("Authorization","Bearer "+token);r.Header.Set("X-GH-Identity-Proof","opaque-proof");r.Header.Set("X-Tenant-ID",uuid.NewString())
  return v.VerifyWorkload(context.Background(),r)
 }
 if _,err=call("GET /identity/restricted/auth/me");err==nil{t.Fatal("disabled actor accepted")}
 if _,err=call("POST /identity/restricted/auth/logout");err!=nil{t.Fatal("logout blocked",err)}
 ctx,err:=call("POST /identity/restricted/auth/login");if err!=nil||school.CheckTenant(ctx,id)==nil{t.Fatal("login policy not propagated",err)}
 active=true
 if _,err=call("GET /identity/restricted/auth/me");err!=nil{t.Fatal("active actor denied",err)}
}
