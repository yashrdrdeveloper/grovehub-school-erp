package registry_test

import (
 "context"
 "crypto/ed25519"
 "crypto/rand"
 "encoding/base64"
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "os"
 "testing"
 "time"

 "github.com/google/uuid"
 "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/postgres"
 core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/registry"
)
func TestSchoolAccessCheck(t *testing.T){
 dsn:=os.Getenv("PROJECT_SCHOOL_TEST_DSN");if dsn==""{t.Skip("requires disposable School database")}
 ctx:=context.Background()
 db,err:=postgres.OpenProjectPool(ctx,dsn);if err!=nil{t.Fatal(err)};defer db.Close()
 var database string
 if err=db.QueryRowContext(ctx,"SELECT current_database()").Scan(&database);err!=nil||database!="school_test"{t.Fatal("disposable database required")}
 if err=registry.Migrate(ctx,db);err!=nil{t.Fatal(err)}
 school,missing:=uuid.NewString(),uuid.NewString()
 if err=registry.RegisterSchool(ctx,db,school,"school-"+school,"Access test school");err!=nil{t.Fatal(err)}
 pub,priv,err:=ed25519.GenerateKey(rand.Reader);if err!=nil{t.Fatal(err)}
 key:=core.WorkloadKeyConfig{KeyID:"access-test",PublicKey:base64.RawURLEncoding.EncodeToString(pub),Issuer:"platform-tenant-access",Subject:"platform-tenant-access",Scopes:[]string{"platform.school.access-check"}}
 raw,_:=json.Marshal([]core.WorkloadKeyConfig{key})
 mux:=http.NewServeMux()
 if err=registry.RegisterHTTP(mux,db,"","","",string(raw));err!=nil{t.Fatal(err)}
 token,err:=core.SignWorkload(key.KeyID,key.Issuer,priv,core.WorkloadClaims{Subject:key.Subject,Audience:"school",Scopes:key.Scopes,ExpiresAt:time.Now().Add(time.Minute)});if err!=nil{t.Fatal(err)}
 var generation string
 check:=func(id,suffix string,auth bool)int{
  r:=httptest.NewRequest("GET","/school/restricted/registry/schools/"+id+"/access-check"+suffix,nil)
  if auth{r.Header.Set("Authorization","Bearer "+token)}
  w:=httptest.NewRecorder();mux.ServeHTTP(w,r);generation=w.Header().Get("X-GH-Session-Generation");return w.Code
 }
 if got:=check(school,"",true);got!=204||generation!="1"{t.Fatal("active school denied or wrong generation",got,generation)}
 if got:=check(missing,"",true);got!=403{t.Fatal("missing school accepted",got)}
 if got:=check(school,"",false);got!=403{t.Fatal("unsigned request accepted",got)}
 if got:=check(school,"?status=ACTIVE",true);got!=403{t.Fatal("query accepted",got)}
 if _,err=db.ExecContext(ctx,"UPDATE school_registry.schools SET status='DISABLED' WHERE id=$1",school);err!=nil{t.Fatal(err)}
 if got:=check(school,"",true);got!=403{t.Fatal("disabled school accepted",got)}
 if _,err=db.ExecContext(ctx,"UPDATE school_registry.schools SET status='ACTIVE' WHERE id=$1",school);err!=nil{t.Fatal(err)}
 if got:=check(school,"",true);got!=204||generation!="2"{t.Fatal("reactivated school denied or wrong generation",got,generation)}
 if _,err=db.ExecContext(ctx,"UPDATE school_registry.schools SET name='Renamed access school' WHERE id=$1",school);err!=nil{t.Fatal(err)}
 if got:=check(school,"",true);got!=204||generation!="2"{t.Fatal("rename changed session generation",got,generation)}
}
