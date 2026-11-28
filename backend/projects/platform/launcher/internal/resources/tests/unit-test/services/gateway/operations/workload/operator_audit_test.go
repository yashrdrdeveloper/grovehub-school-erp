package workload_test

import (
 . "grovehub/launcher/internal/src/services/gateway/operations/workload"
 "context"
 "crypto/ed25519"
 "crypto/rand"
 "encoding/base64"
 "encoding/json"
 "net/http/httptest"
 "testing"
 "time"

 core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
)

func TestOperatorAuditBoundary(t *testing.T){
 pub,priv,err:=ed25519.GenerateKey(rand.Reader);if err!=nil{t.Fatal(err)}
 key:=core.WorkloadKeyConfig{KeyID:"operator-audit-test",PublicKey:base64.RawURLEncoding.EncodeToString(pub),Issuer:"internal-tool-audit",Subject:"internal-tool-audit",Scopes:[]string{OperatorAuditScope}}
 raw,err:=json.Marshal(WorkloadConfig{OperatorAuditKeys:[]core.WorkloadKeyConfig{key}});if err!=nil{t.Fatal(err)}
 verifier,err:=NewPlatformWorkloads(string(raw),nil,nil);if err!=nil{t.Fatal(err)}
 token,err:=core.SignWorkload(key.KeyID,key.Issuer,priv,core.WorkloadClaims{Subject:key.Subject,Audience:"platform",Scopes:key.Scopes,ExpiresAt:time.Now().Add(time.Minute)});if err!=nil{t.Fatal(err)}
 for _,tc:=range []struct{pattern,path,header string;allowed bool}{
  {OperatorAuditRoute,"/audit/restricted/operator-activities","",true},
  {OperatorAuditRoute,"/audit/restricted/operator-activities?tenant=x","",false},
  {OperatorAuditRoute,"/audit/restricted/operator-activities","Cookie",false},
  {OperatorAuditRoute,"/audit/restricted/operator-activities","X-GH-Identity-Proof",false},
  {"POST /audit/restricted/search","/audit/restricted/search","",false},
  {InternalCatalogRoute,"/integrations/restricted/providers","",false},
 }{
  r:=httptest.NewRequest("POST",tc.path,nil);r.Pattern=tc.pattern;r.Header.Set("Authorization","Bearer "+token)
  if tc.header!=""{r.Header.Set(tc.header,"forged")}
  _,err:=verifier.VerifyWorkload(context.Background(),r)
  if (err==nil)!=tc.allowed{t.Errorf("%s %s: %v",tc.pattern,tc.header,err)}
 }
 bad:=key;bad.Scopes=[]string{OperatorAuditScope,"audit.search"}
 if _,err:=NewOperatorAuditVerifier([]core.WorkloadKeyConfig{bad},nil);err==nil{t.Fatal("overbroad key accepted")}
 if _,err:=NewOperatorAuditVerifier([]core.WorkloadKeyConfig{key},[]core.WorkloadKeyConfig{key});err==nil{t.Fatal("shared key accepted")}
}

func TestOperatorAuditReadIsolation(t *testing.T){
 pub,priv,err:=ed25519.GenerateKey(rand.Reader);if err!=nil{t.Fatal(err)}
 key:=core.WorkloadKeyConfig{KeyID:"read-test",PublicKey:base64.RawURLEncoding.EncodeToString(pub),Issuer:"internal-tool-audit-reader",Subject:"internal-tool-audit-reader",Scopes:[]string{OperatorAuditReadScope}}
 raw,_:=json.Marshal(WorkloadConfig{OperatorAuditReadKeys:[]core.WorkloadKeyConfig{key}})
 verifier,err:=NewPlatformWorkloads(string(raw),nil,nil);if err!=nil{t.Fatal(err)}
 token,err:=core.SignWorkload(key.KeyID,key.Issuer,priv,core.WorkloadClaims{Subject:key.Subject,Audience:"platform",Scopes:key.Scopes,ExpiresAt:time.Now().Add(time.Minute)});if err!=nil{t.Fatal(err)}
 for _,path:=range []string{"/audit/restricted/operator-activities/search","/audit/restricted/operator-activities","/audit/restricted/search"}{
  r:=httptest.NewRequest("POST",path,nil);r.Pattern="POST "+path;r.Header.Set("Authorization","Bearer "+token)
  _,err=verifier.VerifyWorkload(context.Background(),r)
  if (err==nil)!=(r.Pattern==OperatorAuditReadRoute){t.Fatal("wrong grant",path,err)}
 }
 if _,err:=NewOperatorAuditReadVerifier([]core.WorkloadKeyConfig{key},[]core.WorkloadKeyConfig{key});err==nil{t.Fatal("shared key accepted")}
}
