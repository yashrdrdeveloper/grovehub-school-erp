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
 bound:=NewWithSchoolBindings(db,nil,map[SchoolBindingKey]SchoolBinding{{EnvironmentCode:env.Code,TemplateID:v1.ID,SchoolID:school}:{Operators:[]string{l1.ID,l2.ID},Validate:func(Invocation)bool{return true},Adapter:func(_ context.Context,call Invocation)(int,error){
  calls++
  if call.RequestID!=a.ID||call.ExecutionID==""||call.Path!="/identity/restricted/accounts/test"||call.SchoolID==nil||*call.SchoolID!=school{t.Error("approved invocation changed")}
  return 200,nil
 }}})
 caps,err:=bound.Capabilities(ctx,l1,a.ID);if err!=nil||!caps.CanExecute||caps.Version!=2{t.Fatal("bound execution capability",caps,err)}
 caps,err=s.Capabilities(ctx,l1,a.ID);if err!=nil||caps.CanExecute{t.Fatal("unbound execution advertised",caps,err)}
 if _,err=bound.Execute(ctx,other,a.ID,ExecutionInput{Version:2});err==nil{t.Fatal("other operator executed request")}
 execution,err:=bound.Execute(ctx,l1,a.ID,ExecutionInput{Version:2});if err!=nil||execution.Status!="SUCCEEDED"||execution.ResponseStatus==nil||*execution.ResponseStatus!=200{t.Fatal("execution failed",execution,err)}
 if _,err=bound.Execute(ctx,l2,a.ID,ExecutionInput{Version:2});!errors.Is(err,ErrConflict)||calls!=1{t.Fatal("duplicate dispatch accepted",err,calls)}
 caps,err=bound.Capabilities(ctx,l1,a.ID);if err!=nil||caps.CanExecute{t.Fatal("duplicate execution advertised",caps,err)}
 saved,err:=bound.Execution(ctx,l2,a.ID);if err!=nil||saved.ID!=execution.ID{t.Fatal("execution history failed",err)}
 first,err:=s.Submit(ctx,l1,req);if err!=nil{t.Fatal(err)}
 second,err:=s.Submit(ctx,l1,req);if err!=nil{t.Fatal(err)}
 bulk:=BulkApprovalInput{Items:[]repo.BulkItem{{ID:first.ID,Version:1},{ID:second.ID,Version:2}}}
 if _,err=s.BulkApprove(ctx,l2,bulk);!errors.Is(err,ErrConflict){t.Fatal("stale bulk accepted",err)}
 check,err:=s.Get(ctx,l1,first.ID);if err!=nil||check.Status!="PENDING"{t.Fatal("bulk rollback failed",err)}
 bulk.Items[1].Version=1
 if _,err=s.BulkApprove(ctx,l1,bulk);!errors.Is(err,ErrForbidden){t.Fatal("L1 bulk accepted")}
 approved,err:=s.BulkApprove(ctx,l2,bulk);if err!=nil||len(approved)!=2||approved[0].Status!="APPROVED"||approved[1].Status!="APPROVED"{t.Fatal("bulk approval failed",err)}
 uncertain,err:=s.Submit(ctx,l2,req);if err!=nil{t.Fatal(err)}
 unknown:=NewWithSchoolBindings(db,nil,map[SchoolBindingKey]SchoolBinding{{EnvironmentCode:env.Code,TemplateID:v1.ID,SchoolID:school}:{Operators:[]string{l2.ID},Validate:func(Invocation)bool{return true},Adapter:func(context.Context,Invocation)(int,error){return 0,errors.New("timeout")}}})
 ex,err:=unknown.Execute(ctx,l2,uncertain.ID,ExecutionInput{Version:uncertain.Version});if err!=nil||ex.Status!="UNKNOWN"{t.Fatal("expected unknown execution",ex,err)}
 evidence:=ReconcileInput{Outcome:"CONFIRMED_FAILED",Evidence:"Verified upstream record"}
 if _,err=s.Reconcile(ctx,l1,uncertain.ID,evidence);!errors.Is(err,ErrForbidden){t.Fatal("L1 reconciled execution")}
 review,err:=s.Reconcile(ctx,l2,uncertain.ID,evidence);if err!=nil||review.ExecutionID!=ex.ID{t.Fatal("reconciliation failed",review,err)}
 if _,err=s.Reconcile(ctx,l2,uncertain.ID,evidence);!errors.Is(err,ErrConflict){t.Fatal("duplicate reconciliation",err)}
 if _,err=unknown.Execute(ctx,l2,uncertain.ID,ExecutionInput{Version:uncertain.Version});!errors.Is(err,ErrConflict){t.Fatal("reconciled execution redispatched",err)}
 p,err:=s.SearchTemplates(ctx,entity.PageRequest{PageNumber:1,RowSize:1});if err!=nil||len(p.Records)!=1||!p.HasMore||p.Records[0].ID!=v2.ID{t.Fatal("template first page",p,err)}
 p,err=s.SearchTemplates(ctx,entity.PageRequest{PageNumber:2,RowSize:1});if err!=nil||len(p.Records)!=1||p.HasMore||p.Records[0].ID!=v1.ID{t.Fatal("template second page",p,err)}
 if _,err=s.SearchApprovals(ctx,l1,entity.PageRequest{},true);!errors.Is(err,ErrForbidden){t.Fatal("L1 accessed inbox search")}
 if _,err=s.SearchEvents(ctx,l1,entity.PageRequest{});!errors.Is(err,ErrForbidden){t.Fatal("L1 accessed event history")}
 ev,err:=s.SearchEvents(ctx,l2,entity.PageRequest{RowSize:1});if err!=nil||len(ev.Records)!=1||!ev.HasMore||ev.Records[0].Action!="EXECUTION_RECONCILED"{t.Fatal("event history",ev,err)}
 restricted:=NewWithSchoolBindings(db,nil,map[SchoolBindingKey]SchoolBinding{{EnvironmentCode:env.Code,TemplateID:v1.ID,SchoolID:school}:{Operators:[]string{l1.ID,l2.ID},Requesters:[]string{l2.ID},Reviewers:[]string{l1.ID},Executors:[]string{l2.ID},Validate:func(Invocation)bool{return true},Adapter:func(context.Context,Invocation)(int,error){return 204,nil}}})
 if _,err=restricted.Submit(ctx,l1,req);!errors.Is(err,ErrForbidden){t.Fatal("unauthorized requester accepted",err)}
 restrictedProposal,err:=restricted.Submit(ctx,l2,req);if err!=nil||restrictedProposal.Status!="PENDING"{t.Fatal("requester without review role auto-approved",restrictedProposal,err)}
 if _,err=restricted.Decide(ctx,l2,restrictedProposal.ID,"APPROVED",DecisionInput{Version:1});!errors.Is(err,ErrForbidden){t.Fatal("unauthorized reviewer accepted",err)}
 if _,err=restricted.BulkApprove(ctx,l2,BulkApprovalInput{Items:[]repo.BulkItem{{ID:restrictedProposal.ID,Version:1}}});!errors.Is(err,ErrForbidden){t.Fatal("unauthorized bulk reviewer accepted",err)}
 mutation:=tpl;mutation.Code="account-patch";mutation.Method="PATCH"
 mutationTemplate,err:=s.CreateTemplate(ctx,l2,mutation);if err!=nil{t.Fatal(err)}
 mutationRequest:=req;mutationRequest.TemplateID=mutationTemplate.ID
 proposal,err:=s.Submit(ctx,l2,mutationRequest);if err!=nil||proposal.Status!="PENDING"{t.Fatal("mutation self-approved",proposal,err)}
 if _,err=s.Decide(ctx,l2,proposal.ID,"APPROVED",DecisionInput{Version:1});!errors.Is(err,ErrForbidden){t.Fatal("mutation self-review accepted",err)}
 if _,err=s.BulkApprove(ctx,l2,BulkApprovalInput{Items:[]repo.BulkItem{{ID:proposal.ID,Version:1}}});!errors.Is(err,ErrForbidden){t.Fatal("bulk mutation self-review accepted",err)}
 var n int
 if err=db.QueryRowContext(ctx,"SELECT count(*) FROM internal_tool.workflow_events WHERE resource_id=$1",a.ID).Scan(&n);err!=nil||n!=4{t.Fatal("audit events missing",n,err)}
}
