package platform

import (
 "context"
 "encoding/json"
 "errors"
 "os"
 "testing"
 "time"

 "github.com/google/uuid"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/services/workflow"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/resources/migrations"
 "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/postgres"
)

// Invoked by the cross-project suite as a separate executable. It owns its
// database and uses the same generated signed clients as the deployed tool.
func TestLiveSchoolLifecycleExecution(t *testing.T) {
 if os.Getenv("LIFECYCLE_EXECUTE")!="true"{t.Skip("cross-project fixture")}
 ctx,cancel:=context.WithTimeout(context.Background(),30*time.Second);defer cancel()
 db,err:=postgres.OpenProjectPool(ctx,os.Getenv("INTERNAL_TOOL_AUTH_TEST_DSN"));if err!=nil{t.Fatal(err)};defer db.Close()
 var name string
 if err=db.QueryRowContext(ctx,"SELECT current_database()").Scan(&name);err!=nil||name!="internal_test"{t.Fatal("disposable Internal Tool database required")}
 if err=migrations.Apply(ctx,db);err!=nil{t.Fatal(err)}
 school:=os.Getenv("LIFECYCLE_SCHOOL_ID")
 origin:=os.Getenv("LIFECYCLE_SCHOOL_ORIGIN")
 dir,err:=NewDirectory(DirectoryConfig{Origin:origin,KeyID:"lifecycle-directory",PrivateKey:os.Getenv("LIFECYCLE_DIRECTORY_PRIVATE"),ReaderID:os.Getenv("LIFECYCLE_DIRECTORY_ID")});if err!=nil{t.Fatal(err)}
 rows,err:=dir(ctx,[]string{school});if err!=nil{t.Fatal(err)}
 current:=rows[school]
 expected,target:=os.Getenv("LIFECYCLE_EXPECTED"),os.Getenv("LIFECYCLE_TARGET")
 if current.Status!=expected||current.Version<1{t.Fatal("unexpected initial registry state",current)}
 l1:=entity.Operator{ID:uuid.NewString(),Persona:entity.L1}
 l2:=entity.Operator{ID:uuid.NewString(),Persona:entity.L2}
 for _,actor:=range []entity.Operator{l1,l2}{
  if _,err=db.ExecContext(ctx,"INSERT INTO internal_tool.operators(id,email,persona,password_salt,password_hash) VALUES($1,$2,$3,$4,$5)",actor.ID,actor.ID+"@example.test",string(actor.Persona),make([]byte,16),make([]byte,32));err!=nil{t.Fatal(err)}
 }
 base:=workflow.New(db)
 env,err:=base.CreateEnvironment(ctx,l2,workflow.EnvironmentInput{Code:"lifecycle-"+uuid.NewString(),Name:"Lifecycle test"});if err!=nil{t.Fatal(err)}
 tpl,err:=base.CreateTemplate(ctx,l2,workflow.TemplateInput{Code:"lifecycle-"+uuid.NewString(),Module:"school",Method:"PATCH",Endpoint:"/school/restricted/registry/schools/${schoolId}/status",Scope:"SCHOOL",BodySchema:map[string]entity.Field{"schoolId":{Type:"string",Required:true},"expectedStatus":{Type:"string",Required:true},"status":{Type:"string",Required:true},"expectedVersion":{Type:"integer",Required:true}}});if err!=nil{t.Fatal(err)}
 bindings,err:=StatusBindings(SchoolConfig{EnvironmentCode:env.Code,TemplateID:tpl.ID,SchoolID:school,Origin:origin,KeyID:"lifecycle-status",PrivateKey:os.Getenv("LIFECYCLE_STATUS_PRIVATE"),OperatorIDs:l1.ID+","+l2.ID});if err!=nil{t.Fatal(err)}
 svc:=workflow.NewWithSchoolBindings(db,nil,bindings).WithDirectory(dir)
 raw,err:=json.Marshal(map[string]any{"schoolId":school,"expectedStatus":expected,"status":target,"expectedVersion":current.Version});if err!=nil{t.Fatal(err)}
 input:=workflow.SubmitInput{TemplateID:tpl.ID,EnvironmentID:env.ID,SchoolID:&school,Input:raw,Reason:"Cross-project lifecycle verification"}
 request,err:=svc.Submit(ctx,l1,input);if err!=nil||request.Status!="PENDING"{t.Fatal("submit",request,err)}
 if _,err=svc.Execute(ctx,l1,request.ID,workflow.ExecutionInput{Version:request.Version});err==nil{t.Fatal("unapproved change executed")}
 if _,err=svc.Decide(ctx,l1,request.ID,"APPROVED",workflow.DecisionInput{Version:request.Version});!errors.Is(err,workflow.ErrForbidden){t.Fatal("L1 approved own request",err)}
 approved,err:=svc.Decide(ctx,l2,request.ID,"APPROVED",workflow.DecisionInput{Version:request.Version});if err!=nil||approved.Status!="APPROVED"{t.Fatal("approve",approved,err)}
 rows,err=dir(ctx,[]string{school});if err!=nil||rows[school].Status!=expected||rows[school].Version!=current.Version{t.Fatal("approval alone mutated registry",rows,err)}
 ex,err:=svc.Execute(ctx,l1,request.ID,workflow.ExecutionInput{Version:approved.Version});if err!=nil||ex.Status!="SUCCEEDED"||ex.ResponseStatus==nil||*ex.ResponseStatus!=204{t.Fatal("execute",ex,err)}
 rows,err=dir(ctx,[]string{school});if err!=nil||rows[school].Status!=target||rows[school].Version!=current.Version+1{t.Fatal("committed registry state",rows,err)}
 if _,err=svc.Execute(ctx,l1,request.ID,workflow.ExecutionInput{Version:approved.Version});!errors.Is(err,workflow.ErrConflict){t.Fatal("duplicate execution accepted",err)}
}
