package registry_test

import (
 "bytes"
 "context"
 "errors"
 sdk "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/internal/generated/audit"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/registry/internal/app/services/delivery"
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
 coredatabase "grovehub/libs/core/database"
 core "grovehub/libs/core/gateway"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/registry"
)
func TestSchoolStatus(t *testing.T){
 dsn:=os.Getenv("PROJECT_SCHOOL_TEST_DSN");if dsn==""{t.Skip("requires disposable School database")}
 ctx:=context.Background()
 db,err:=coredatabase.OpenProjectPool(ctx,dsn);if err!=nil{t.Fatal(err)};defer db.Close()
 var database string
 if err=db.QueryRowContext(ctx,"SELECT current_database()").Scan(&database);err!=nil||database!="school_test"{t.Fatal("disposable database required")}
 if err=registry.Migrate(ctx,db);err!=nil{t.Fatal(err)}
 school,other:=uuid.NewString(),uuid.NewString()
 for _,id:=range []string{school,other}{if err=registry.RegisterSchool(ctx,db,id,"school-"+id,"Status test school");err!=nil{t.Fatal(err)}}
 pub,priv,err:=ed25519.GenerateKey(rand.Reader);if err!=nil{t.Fatal(err)}
 key:=core.WorkloadKeyConfig{KeyID:"status-test",PublicKey:base64.RawURLEncoding.EncodeToString(pub),Issuer:"internal-tool-school-status",Subject:"internal-tool-school-status:"+school,Scopes:[]string{"internal.school.registry.status"}}
 raw,_:=json.Marshal([]core.WorkloadKeyConfig{key})
 mux:=http.NewServeMux()
 if err=registry.RegisterHTTP(mux,db,"","",string(raw));err!=nil{t.Fatal(err)}
 token,err:=core.SignWorkload(key.KeyID,key.Issuer,priv,core.WorkloadClaims{Subject:key.Subject,Audience:"school",Scopes:key.Scopes,ExpiresAt:time.Now().Add(time.Minute)});if err!=nil{t.Fatal(err)}
 correlationID:=uuid.NewString()
 call:=func(id,execution,expected,status string,version int64,auth bool)int{
  body,_:=json.Marshal(map[string]any{"expectedStatus":expected,"status":status,"expectedVersion":version,"correlationId":correlationID})
  r:=httptest.NewRequest("PATCH","/school/restricted/registry/schools/"+id+"/status",bytes.NewReader(body))
  r.Header.Set("Idempotency-Key",execution)
  if auth{r.Header.Set("Authorization","Bearer "+token)}
  w:=httptest.NewRecorder();mux.ServeHTTP(w,r);return w.Code
 }
 execution:=uuid.NewString()
 for _,tc:=range []struct{id,execution,expected,status string;version int64;auth bool;want int}{
  {other,execution,"ACTIVE","DISABLED",1,true,403},
  {school,execution,"ACTIVE","DISABLED",1,false,403},
  {school,execution,"ACTIVE","ACTIVE",1,true,400},
  {school,execution,"ACTIVE","DISABLED",2,true,409},
  {school,execution,"ACTIVE","DISABLED",1,true,204},
  {school,execution,"ACTIVE","DISABLED",1,true,204},
  {school,execution,"DISABLED","ACTIVE",2,true,409},
  {school,uuid.NewString(),"DISABLED","ACTIVE",2,true,204},
  {school,execution,"ACTIVE","DISABLED",1,true,204},
  {school,uuid.NewString(),"ACTIVE","DISABLED",1,true,409},
 }{if got:=call(tc.id,tc.execution,tc.expected,tc.status,tc.version,tc.auth);got!=tc.want{t.Fatal("unexpected status result",tc.want,got)}}
 var status string;var version int64;var count int
 if err=db.QueryRowContext(ctx,"SELECT status,version FROM school_registry.schools WHERE id=$1",school).Scan(&status,&version);err!=nil||status!="ACTIVE"||version!=3{t.Fatal("incorrect status/version",status,version,err)}
 if err=db.QueryRowContext(ctx,"SELECT count(*) FROM school_registry.status_events WHERE school_id=$1 AND correlation_id=$2",school,correlationID).Scan(&count);err!=nil||count!=2{t.Fatal("incorrect receipts",count,err)}
 if _,err=db.ExecContext(ctx,"DELETE FROM school_registry.status_events WHERE execution_id=$1",execution);err==nil{t.Fatal("receipt deleted")}
 if err=db.QueryRowContext(ctx,"SELECT count(*) FROM school_registry.status_outbox o JOIN school_registry.status_events e USING(execution_id) WHERE e.school_id=$1",school).Scan(&count);err!=nil||count!=2{t.Fatal("status outbox count",count,err)}
 calls:=0;var first string
 worker:=delivery.NewStatus(db,school,func(_ context.Context,a sdk.SchoolStatusReceipt)(int,error){
  calls++
  if a.SchoolId!=school||a.CorrelationId==nil||*a.CorrelationId!=correlationID{t.Error("receipt identity mismatch")}
  if (a.PreviousVersion==1&&a.Operation!="DISABLE")||(a.PreviousVersion==2&&a.Operation!="ACTIVATE"){t.Error("incorrect transition")}
  if calls==1{first=a.ExecutionId;return 0,errors.New("temporary failure")}
  if calls==2&&a.ExecutionId!=first{t.Error("retry identity changed")}
  return 204,nil
 })
 if worked,err:=worker.Step(ctx);!worked||err!=nil{t.Fatal("first delivery",worked,err)}
 if _,err=db.ExecContext(ctx,"UPDATE school_registry.status_outbox SET next_attempt_at='2000-01-01' WHERE execution_id=$1",first);err!=nil{t.Fatal(err)}
 for n:=0;n<2;n++{if worked,err:=worker.Step(ctx);!worked||err!=nil{t.Fatal("delivery",worked,err)}}
 if worked,err:=worker.Step(ctx);worked||err!=nil{t.Fatal("unexpected replay",worked,err)}
 if err=db.QueryRowContext(ctx,"SELECT count(*) FROM school_registry.status_outbox o JOIN school_registry.status_events e USING(execution_id) WHERE e.school_id=$1 AND o.delivered_at IS NOT NULL",school).Scan(&count);err!=nil||count!=2{t.Fatal("delivered count",count,err)}
}
