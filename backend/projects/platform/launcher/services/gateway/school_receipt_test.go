package gateway

import(
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
func TestSchoolReceiptWorkload(t *testing.T){
 pub,priv,err:=ed25519.GenerateKey(rand.Reader);if err!=nil{t.Fatal(err)}
 k:=core.WorkloadKeyConfig{KeyID:"receipt-test",PublicKey:base64.RawURLEncoding.EncodeToString(pub),Issuer:receiptIssuer,Subject:receiptIssuer+":"+uuid.NewString(),Scopes:[]string{receiptScope}}
 raw,_:=json.Marshal(workloadConfig{SchoolReceiptKeys:[]core.WorkloadKeyConfig{k}})
 v,err:=NewPlatformWorkloads(string(raw),nil,nil);if err!=nil{t.Fatal(err)}
 token,err:=core.SignWorkload(k.KeyID,k.Issuer,priv,core.WorkloadClaims{Subject:k.Subject,Audience:"platform",Scopes:k.Scopes,ExpiresAt:time.Now().Add(time.Minute)});if err!=nil{t.Fatal(err)}
 for _,path:=range []string{"/audit/restricted/school-rename-receipts","/audit/restricted/operator-activities","/audit/restricted/search"}{
  r:=httptest.NewRequest("POST",path,nil);r.Pattern="POST "+path;r.Header.Set("Authorization","Bearer "+token)
  _,err=v.VerifyWorkload(context.Background(),r)
  if (err==nil)!=(r.Pattern==receiptRoute){t.Fatal("wrong grant",path,err)}
 }
 if _,err:=newReceiptVerifier([]core.WorkloadKeyConfig{k},[]core.WorkloadKeyConfig{k});err==nil{t.Fatal("shared credential accepted")}
}

func TestSchoolReceiptReaderWorkload(t *testing.T){
 pub,priv,err:=ed25519.GenerateKey(rand.Reader);if err!=nil{t.Fatal(err)}
 k:=core.WorkloadKeyConfig{KeyID:"receipt-reader-test",PublicKey:base64.RawURLEncoding.EncodeToString(pub),Issuer:"internal-tool-school-receipt-reader",Subject:"internal-tool-school-receipt-reader",Scopes:[]string{receiptReadScope}}
 raw,_:=json.Marshal(workloadConfig{SchoolReceiptReadKeys:[]core.WorkloadKeyConfig{k}})
 v,err:=NewPlatformWorkloads(string(raw),nil,nil);if err!=nil{t.Fatal(err)}
 token,err:=core.SignWorkload(k.KeyID,k.Issuer,priv,core.WorkloadClaims{Subject:k.Subject,Audience:"platform",Scopes:k.Scopes,ExpiresAt:time.Now().Add(time.Minute)});if err!=nil{t.Fatal(err)}
 for _,path:=range []string{"/audit/restricted/school-mutation-receipts/search","/audit/restricted/school-rename-receipts","/audit/restricted/operator-activities/search","/audit/restricted/search"}{
  r:=httptest.NewRequest("POST",path,nil);r.Pattern="POST "+path;r.Header.Set("Authorization","Bearer "+token)
  _,err=v.VerifyWorkload(context.Background(),r)
  if (err==nil)!=(r.Pattern==receiptReadRoute){t.Fatal("wrong reader grant",path,err)}
 }
 if _,err:=newReceiptReadVerifier([]core.WorkloadKeyConfig{k},[]core.WorkloadKeyConfig{k});err==nil{t.Fatal("reused key accepted")}
}
