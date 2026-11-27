package delivery

import (
 "crypto/ed25519"
 "crypto/rand"
 "encoding/base64"
 "encoding/json"
 "testing"

 "github.com/google/uuid"
 core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
)
func TestConfigure(t *testing.T){
 makeConfig:=func()Config{
  _,key,err:=ed25519.GenerateKey(rand.Reader);if err!=nil{t.Fatal(err)}
  return Config{Origin:"http://127.0.0.1:8080",SchoolID:uuid.NewString(),KeyID:uuid.NewString(),PrivateKey:base64.RawURLEncoding.EncodeToString(key)}
 }
 a,b:=makeConfig(),makeConfig()
 encode:=func(v any)string{b,err:=json.Marshal(v);if err!=nil{t.Fatal(err)};return string(b)}
 status:=b;status.SchoolID=a.SchoolID;status.Operation="status"
 if configs,err:=Configure(encode([]Config{a,status}),Config{},"","","","");err!=nil||len(configs)!=2{t.Fatal("separate operation publishers rejected",err)}
 duplicate:=status;duplicate.KeyID=uuid.NewString()
 _,private,err:=ed25519.GenerateKey(rand.Reader);if err!=nil{t.Fatal(err)}
 duplicate.PrivateKey=base64.RawURLEncoding.EncodeToString(private)
 if _,err:=Configure(encode([]Config{a,status,duplicate}),Config{},"","","","");err==nil{t.Fatal("duplicate lifecycle publisher accepted")}
 raw:=encode([]Config{a,b})
 if configs,err:=Configure(raw,Config{},"","","","");err!=nil||len(configs)!=2{t.Fatal("two schools rejected",err)}
 if _,err:=Configure(raw,a,"","","","");err==nil{t.Fatal("mixed configuration")}
 for _,bad:=range []Config{a,func()Config{c:=b;c.KeyID=a.KeyID;return c}(),func()Config{c:=b;c.PrivateKey=a.PrivateKey;return c}(),func()Config{c:=b;c.SchoolID=a.SchoolID;return c}()}{
  if _,err:=Configure(encode([]Config{a,bad}),Config{},"","","","");err==nil{t.Fatal("reused credential or school accepted")}
 }
 private,_=base64.RawURLEncoding.DecodeString(a.PrivateKey)
 public:=base64.RawURLEncoding.EncodeToString(private[32:])
 readers:=encode([]ReaderGrant{{Key:core.WorkloadKeyConfig{KeyID:"reader",PublicKey:public}}})
 if _,err:=Configure(raw,Config{},"","",readers,"");err==nil{t.Fatal("reader key reused")}
 if _,err:=Configure(raw,Config{},"other",a.PrivateKey,"","");err==nil{t.Fatal("outbound key reused")}
 if _,err:=Configure(raw,Config{},"","","",encode([]core.WorkloadKeyConfig{{KeyID:a.KeyID}}));err==nil{t.Fatal("rename key ID reused")}
}
