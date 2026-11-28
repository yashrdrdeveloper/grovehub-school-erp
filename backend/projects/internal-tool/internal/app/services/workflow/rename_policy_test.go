package workflow

import (
 "context"
 "errors"
 "testing"
 "github.com/google/uuid"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
)
func TestRenameExecutionPolicy(t *testing.T){
 school,revision,requester,reviewer:=uuid.NewString(),uuid.NewString(),uuid.NewString(),uuid.NewString()
 binding:=SchoolBinding{Mutation:true,Operators:[]string{requester,reviewer},Validate:func(Invocation)bool{return true},Adapter:func(context.Context,Invocation)(int,error){return 204,nil}}
 s:=NewWithSchoolBindings(nil,nil,map[SchoolBindingKey]SchoolBinding{{"dev",revision,school}:binding})
 tpl:=entity.Template{ID:revision,Scope:"SCHOOL",Method:"PATCH"}
 a:=entity.Approval{SchoolID:&school,RequesterID:requester,ReviewerID:&reviewer}
 actor:=entity.Operator{ID:requester,Persona:entity.L2}
 call:=Invocation{SchoolID:&school,Method:"PATCH"}
 if _,err:=s.executionAdapter("dev",tpl,a,actor,call);!errors.Is(err,ErrForbidden){t.Fatal("rename without directory accepted",err)}
 s.WithDirectory(func(context.Context,[]string)(map[string]SchoolRecord,error){return map[string]SchoolRecord{school:{ID:school,Status:"ACTIVE"}},nil})
 if _,err:=s.executionAdapter("dev",tpl,a,actor,call);err!=nil{t.Fatal("authorized rename rejected",err)}
 a.ReviewerID=&requester
 if _,err:=s.executionAdapter("dev",tpl,a,actor,call);!errors.Is(err,ErrForbidden){t.Fatal("self-approved rename accepted",err)}
 a.ReviewerID=&reviewer;call.Method="DELETE"
 if _,err:=s.executionAdapter("dev",tpl,a,actor,call);!errors.Is(err,ErrForbidden){t.Fatal("unbound mutation accepted",err)}
}
