package workload_test

import (
 "grovehub/launcher/internal/src/services/gateway/operations/workload/strategies/internalcatalog"
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

func TestInternalCatalogBoundary(t *testing.T){
 pub,priv,err:=ed25519.GenerateKey(rand.Reader);if err!=nil{t.Fatal(err)}
 key:=core.WorkloadKeyConfig{KeyID:"internal-test",PublicKey:base64.RawURLEncoding.EncodeToString(pub),Issuer:"internal-tool",Subject:"internal-tool",Scopes:[]string{internalcatalog.InternalCatalogScope}}
 raw,err:=json.Marshal(WorkloadConfig{InternalToolKeys:[]core.WorkloadKeyConfig{key}});if err!=nil{t.Fatal(err)}
 verifier,err:=NewPlatformWorkloads(string(raw),nil,nil);if err!=nil{t.Fatal(err)}
 token,err:=core.SignWorkload(key.KeyID,key.Issuer,priv,core.WorkloadClaims{Subject:key.Subject,Audience:"platform",Scopes:key.Scopes,ExpiresAt:time.Now().Add(time.Minute)});if err!=nil{t.Fatal(err)}
 for _,tc:=range []struct{pattern,path,header string;allowed bool}{
  {internalcatalog.InternalCatalogRoute,"/integrations/restricted/providers","",true},
  {internalcatalog.InternalCatalogRoute,"/integrations/restricted/providers?school=x","",false},
  {internalcatalog.InternalCatalogRoute,"/integrations/restricted/providers","Cookie",false},
  {internalcatalog.InternalCatalogRoute,"/integrations/restricted/providers","X-GH-Identity-Proof",false},
  {"GET /integrations/restricted/connections","/integrations/restricted/connections","",false},
  {"POST /foundation/restricted/persons","/foundation/restricted/persons","",false},
 }{
  r:=httptest.NewRequest("GET",tc.path,nil);r.Pattern=tc.pattern;r.Header.Set("Authorization","Bearer "+token)
  if tc.header!=""{r.Header.Set(tc.header,"forged")}
  _,err:=verifier.VerifyWorkload(context.Background(),r)
  if (err==nil)!=tc.allowed{t.Errorf("%s %s: %v",tc.pattern,tc.header,err)}
 }
 bad:=key;bad.Scopes=[]string{internalcatalog.InternalCatalogScope,"identity.accounts.create"}
 if _,err:=internalcatalog.New([]core.WorkloadKeyConfig{bad},nil);err==nil{t.Fatal("overbroad key accepted")}
 if _,err:=internalcatalog.New([]core.WorkloadKeyConfig{key},[]core.WorkloadKeyConfig{key});err==nil{t.Fatal("shared key accepted")}
}
