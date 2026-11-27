package workflow

import (
 "context"
 "encoding/json"
 "errors"
 "os"
 "testing"

 "github.com/google/uuid"
 repo "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/repository/workflow"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/resources/migrations"
 "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/postgres"
)

func TestApprovalLifecycle(t *testing.T){
 dsn:=os.Getenv("INTERNAL_TOOL_AUTH_TEST_DSN");if dsn==""{t.Skip("requires disposable database")}
 ctx:=context.Background();db,err:=postgres.OpenProjectPool(ctx,dsn);if err!=nil{t.Fatal(err)};defer db.Close()
 if err=migrations.Apply(ctx,db);err!=nil{t.Fatal(err)}
 l1:=entity.Operator{ID:uuid.NewString(),Email:"workflow-l1@example.com",Persona:entity.L1}
 l2:=entity.Operator{ID:uuid.NewString(),Email:"workflow-l2@example.com",Persona:entity.L2}
 other:=entity.Operator{ID:uuid.NewString(),Email:"workflow-other@example.com",Persona:entity.L1}
 for _,op:=range []entity.Operator{l1,l2,other}{if _,err=db.ExecContext(ctx,"INSERT INTO internal_tool.operators(id,email,persona,password_salt,password_hash) VALUES($1,$2,$3,$4,$5)",op.ID,op.Email,string(op.Persona),make([]byte,16),make([]byte,32));err!=nil{t.Fatal(err)}}
 s:=New(db)
 input:=EnvironmentInput{Code:"test",Name:"Test"}
 if _,err=s.CreateEnvironment(ctx,l1,input);!errors.Is(err,ErrForbidden){t.Fatal("L1 registered environment")}
 env,err:=s.CreateEnvironment(ctx,l2,input);if err!=nil{t.Fatal(err)}
 tpl:=TemplateInput{Code:"account-get",Module:"identity",Method:"GET",Endpoint:"/identity/restricted/accounts/${accountId}",Scope:"SCHOOL",BodySchema:map[string]entity.Field{"accountId":{Type:"string",Required:true}}}
 v1,err:=s.CreateTemplate(ctx,l2,tpl);if err!=nil{t.Fatal(err)}
 v2,err:=s.CreateTemplate(ctx,l2,tpl);if err!=nil||v2.Version!=2||v1.ID==v2.ID{t.Fatal("immutable revision failed",err)}
 school:=uuid.NewString()
 req:=SubmitInput{TemplateID:v1.ID,EnvironmentID:env.ID,SchoolID:&school,Input:json.RawMessage(`{"accountId":"test"}`),Reason:"Review access"}
 a,err:=s.Submit(ctx,l1,req);if err!=nil||a.Status!="PENDING"||a.Version!=1{t.Fatal("submission failed",a,err)}
 if _,err=s.List(ctx,l1,true);!errors.Is(err,ErrForbidden){t.Fatal("L1 inbox access")}
 if _,err=s.Get(ctx,other,a.ID);err==nil{t.Fatal("other L1 read request")}
 if _,err=s.Decide(ctx,l1,a.ID,"APPROVED",DecisionInput{Version:1});!errors.Is(err,ErrForbidden){t.Fatal("L1 approved request")}
 if _,err=s.Decide(ctx,other,a.ID,"CANCELLED",DecisionInput{Version:1});!errors.Is(err,ErrForbidden){t.Fatal("nonowner cancelled")}
 b,err:=s.Decide(ctx,l2,a.ID,"APPROVED",DecisionInput{Version:1});if err!=nil||b.Status!="APPROVED"||b.Version!=2||b.ReviewerID==nil||*b.ReviewerID!=l2.ID||b.TemplateID!=v1.ID{t.Fatal("approval failed",b,err)}
 if _,err=s.Decide(ctx,l2,a.ID,"REJECTED",DecisionInput{Version:1,Comment:"late"});!errors.Is(err,ErrConflict){t.Fatal("stale decision accepted")}
 auto,err:=s.Submit(ctx,l2,req);if err!=nil||auto.Status!="APPROVED"||auto.ReviewerID==nil||*auto.ReviewerID!=l2.ID{t.Fatal("L2 self approval failed",auto,err)}
 pending,err:=s.Submit(ctx,l1,req);if err!=nil{t.Fatal(err)}
 cancelled,err:=s.Decide(ctx,l1,pending.ID,"CANCELLED",DecisionInput{Version:1});if err!=nil||cancelled.Status!="CANCELLED"{t.Fatal("cancel failed",err)}
 if _,err=s.Execute(ctx,l1,a.ID,ExecutionInput{Version:2});!errors.Is(err,ErrExecutionUnavailable){t.Fatal("unbound execution not denied",err)}
 calls:=0
 bound:=NewWithBindings(db,map[BindingKey]Adapter{{EnvironmentCode:env.Code,TemplateID:v1.ID}:func(_ context.Context,call Invocation)(int,error){
  calls++
  if call.RequestID!=a.ID||call.ExecutionID==""||call.Path!="/identity/restricted/accounts/test"||call.SchoolID==nil||*call.SchoolID!=school{t.Error("approved invocation changed")}
  return 200,nil
 }})
 if _,err=bound.Execute(ctx,other,a.ID,ExecutionInput{Version:2});err==nil{t.Fatal("other operator executed request")}
 execution,err:=bound.Execute(ctx,l1,a.ID,ExecutionInput{Version:2});if err!=nil||execution.Status!="SUCCEEDED"||execution.ResponseStatus==nil||*execution.ResponseStatus!=200{t.Fatal("execution failed",execution,err)}
 if _,err=bound.Execute(ctx,l2,a.ID,ExecutionInput{Version:2});!errors.Is(err,ErrConflict)||calls!=1{t.Fatal("duplicate dispatch accepted",err,calls)}
 saved,err:=bound.Execution(ctx,l2,a.ID);if err!=nil||saved.ID!=execution.ID{t.Fatal("execution history failed",err)}
 first,err:=s.Submit(ctx,l1,req);if err!=nil{t.Fatal(err)}
 second,err:=s.Submit(ctx,l1,req);if err!=nil{t.Fatal(err)}
 bulk:=BulkApprovalInput{Items:[]repo.BulkItem{{ID:first.ID,Version:1},{ID:second.ID,Version:2}}}
 if _,err=s.BulkApprove(ctx,l2,bulk);!errors.Is(err,ErrConflict){t.Fatal("stale bulk accepted",err)}
 check,err:=s.Get(ctx,l1,first.ID);if err!=nil||check.Status!="PENDING"{t.Fatal("bulk rollback failed",err)}
 bulk.Items[1].Version=1
 if _,err=s.BulkApprove(ctx,l1,bulk);!errors.Is(err,ErrForbidden){t.Fatal("L1 bulk accepted")}
 approved,err:=s.BulkApprove(ctx,l2,bulk);if err!=nil||len(approved)!=2||approved[0].Status!="APPROVED"||approved[1].Status!="APPROVED"{t.Fatal("bulk approval failed",err)}
 var n int
 if err=db.QueryRowContext(ctx,"SELECT count(*) FROM internal_tool.workflow_events WHERE resource_id=$1",a.ID).Scan(&n);err!=nil||n!=4{t.Fatal("audit events missing",n,err)}
}
