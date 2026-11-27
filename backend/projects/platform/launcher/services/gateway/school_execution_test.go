package gateway

import (
 "context"
 "crypto/ed25519"
 "crypto/rand"
 "encoding/base64"
 "encoding/json"
 "net/http/httptest"
 "testing"
 "time"

 "github.com/google/uuid"
 core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
)
func TestSchoolExecutionBoundary(t *testing.T){
 school,other:=uuid.NewString(),uuid.NewString()
 pub,priv,err:=ed25519.GenerateKey(rand.Reader);if err!=nil{t.Fatal(err)}
 key:=core.WorkloadKeyConfig{KeyID:"school-check",PublicKey:base64.RawURLEncoding.EncodeToString(pub),Issuer:schoolExecutionIssuer,Subject:schoolSubjectPrefix+school,Scopes:[]string{schoolAuditCheckScope}}
 raw,_:=json.Marshal(workloadConfig{SchoolExecutionKeys:[]core.WorkloadKeyConfig{key}})
 verifier,err:=NewPlatformWorkloads(string(raw),nil,nil);if err!=nil{t.Fatal(err)}
 token,err:=core.SignWorkload(key.KeyID,key.Issuer,priv,core.WorkloadClaims{Subject:key.Subject,Audience:"platform",Scopes:key.Scopes,ExpiresAt:time.Now().Add(time.Minute)});if err!=nil{t.Fatal(err)}
 for _,tc:=range []struct{school,method,suffix,header string;allow bool}{
  {school,"GET","","",true},{other,"GET","","",false},{school,"GET","?tenant="+other,"",false},{school,"POST","","",false},{school,"GET","","Cookie",false},{school,"GET","","X-GH-Identity-Proof",false},
 }{
  r:=httptest.NewRequest(tc.method,"/audit/restricted/schools/"+tc.school+"/activity-check"+tc.suffix,nil)
  r.Pattern=schoolAuditCheckRoute;r.Header.Set("Authorization","Bearer "+token)
  if tc.header!=""{r.Header.Set(tc.header,"forged")}
  _,err=verifier.VerifyWorkload(context.Background(),r)
  if (err==nil)!=tc.allow{t.Fatal("unexpected authorization",tc,err)}
 }
 bad:=key;bad.Subject=schoolSubjectPrefix+other
 if _,err:=newSchoolExecutionVerifier([]core.WorkloadKeyConfig{key,bad},nil);err==nil{t.Fatal("shared school key accepted")}
 if _,err:=newSchoolExecutionVerifier([]core.WorkloadKeyConfig{key},[]core.WorkloadKeyConfig{key});err==nil{t.Fatal("shared workload accepted")}
}
