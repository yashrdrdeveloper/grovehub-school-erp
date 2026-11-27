package auditdelivery

import (
 "context"
 "errors"
 "os"
 "testing"

 "github.com/google/uuid"
 "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/postgres"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/client/audit"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/resources/migrations"
)

func TestDurableRelay(t *testing.T){
 dsn:=os.Getenv("INTERNAL_TOOL_AUTH_TEST_DSN");if dsn==""{t.Skip("requires disposable database")}
 ctx:=context.Background();db,err:=postgres.OpenProjectPool(ctx,dsn);if err!=nil{t.Fatal(err)};defer db.Close()
 if err=migrations.Apply(ctx,db);err!=nil{t.Fatal(err)}
 // Use a dedicated database in CI; mark earlier test rows delivered before
 // creating a new occurrence so the worker selects exactly our test row.
 if _,err=db.ExecContext(ctx,"UPDATE internal_tool.audit_outbox SET delivered_at=now() WHERE blocked_at IS NULL");err!=nil{t.Fatal(err)}
 actor,resource:=uuid.NewString(),uuid.NewString()
 if _,err=db.ExecContext(ctx,"INSERT INTO internal_tool.operators(id,email,persona,password_salt,password_hash) VALUES($1,$2,'L2',$3,$4)",actor,actor+"@example.com",make([]byte,16),make([]byte,32));err!=nil{t.Fatal(err)}
 tx,err:=db.BeginTx(ctx,nil);if err!=nil{t.Fatal(err)}
 if _,err=tx.ExecContext(ctx,"INSERT INTO internal_tool.workflow_events(actor_id,resource_id,action) VALUES($1,$2,'REQUEST_SUBMITTED')",actor,resource);err!=nil{t.Fatal(err)}
 if err=tx.Rollback();err!=nil{t.Fatal(err)}
 var count int
 if err=db.QueryRowContext(ctx,"SELECT count(*) FROM internal_tool.workflow_events WHERE resource_id=$1",resource).Scan(&count);err!=nil||count!=0{t.Fatal("rollback",count,err)}
 var id int64
 if err=db.QueryRowContext(ctx,"INSERT INTO internal_tool.workflow_events(actor_id,resource_id,action) VALUES($1,$2,'REQUEST_SUBMITTED') RETURNING id",actor,resource).Scan(&id);err!=nil{t.Fatal(err)}
 var eventID string
 calls:=0
 w:=New(db,func(_ context.Context,a audit.Activity)(int,error){
  calls++
  if a.ActorID!=actor||a.ResourceID!=resource||a.Operation!="REQUEST_SUBMITTED"{t.Error("incorrect event")}
  if calls==1{eventID=a.EventID;return 0,errors.New("network failure")}
  if a.EventID!=eventID{t.Error("retry changed occurrence ID")}
  return 204,nil
 })
 if worked,err:=w.Step(ctx);!worked||err!=nil{t.Fatal("first attempt",worked,err)}
 if worked,err:=w.Step(ctx);worked||err!=nil{t.Fatal("retry not delayed",worked,err)}
 if _,err=db.ExecContext(ctx,"UPDATE internal_tool.audit_outbox SET next_attempt_at=now() WHERE workflow_event_id=$1",id);err!=nil{t.Fatal(err)}
 if worked,err:=w.Step(ctx);!worked||err!=nil{t.Fatal("second attempt",worked,err)}
 if worked,err:=w.Step(ctx);worked||err!=nil||calls!=2{t.Fatal("delivered row retried",worked,err,calls)}
 var attempts int;var delivered bool
 if err=db.QueryRowContext(ctx,"SELECT attempts,delivered_at IS NOT NULL FROM internal_tool.audit_outbox WHERE workflow_event_id=$1",id).Scan(&attempts,&delivered);err!=nil||attempts!=2||!delivered{t.Fatal("delivery receipt",attempts,delivered,err)}
}
