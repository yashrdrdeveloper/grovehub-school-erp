package workflow

import (
 "context"
 "errors"
 "testing"

 "github.com/google/uuid"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
)
func TestSchoolExecutionPolicy(t *testing.T){
 school,other,template,requester,reviewer:=uuid.NewString(),uuid.NewString(),uuid.NewString(),uuid.NewString(),uuid.NewString()
 key:=SchoolBindingKey{"dev",template,school}
 adapter:=func(context.Context,Invocation)(int,error){return 204,nil}
 s:=NewWithSchoolBindings(nil,map[BindingKey]Adapter{{"dev",template}:adapter},map[SchoolBindingKey]SchoolBinding{key:{Operators:[]string{requester,reviewer},Validate:func(c Invocation)bool{return c.Path=="/allowed"},Adapter:adapter}})
 tpl:=entity.Template{ID:template,Scope:"SCHOOL",Method:"GET"}
 a:=entity.Approval{SchoolID:&school,RequesterID:requester,ReviewerID:&reviewer}
 actor:=entity.Operator{ID:requester,Persona:entity.L1}
 call:=Invocation{SchoolID:&school,Method:"GET",Path:"/allowed"}
 if _,err:=s.executionAdapter("dev",tpl,a,actor,call);err!=nil{t.Fatal(err)}
 if _,err:=s.executionAdapter("other",tpl,a,actor,call);!errors.Is(err,ErrExecutionUnavailable){t.Fatal("other environment",err)}
 a.SchoolID=&other
 if _,err:=s.executionAdapter("dev",tpl,a,actor,call);!errors.Is(err,ErrExecutionUnavailable){t.Fatal("other school",err)}
 a.SchoolID=&school
 actor.ID=other
 if _,err:=s.executionAdapter("dev",tpl,a,actor,call);!errors.Is(err,ErrForbidden){t.Fatal("other operator",err)}
 actor.ID=requester;a.ReviewerID=&other
 if _,err:=s.executionAdapter("dev",tpl,a,actor,call);!errors.Is(err,ErrForbidden){t.Fatal("other reviewer",err)}
 a.ReviewerID=&reviewer;call.Path="/wrong"
 if _,err:=s.executionAdapter("dev",tpl,a,actor,call);!errors.Is(err,ErrForbidden){t.Fatal("wrong path",err)}
 call.Path="/allowed";tpl.Method="PATCH"
 if _,err:=s.executionAdapter("dev",tpl,a,actor,call);!errors.Is(err,ErrForbidden){t.Fatal("mutation accepted",err)}
 tpl.Method="GET"
 globalOnly:=NewWithBindings(nil,map[BindingKey]Adapter{{"dev",template}:adapter})
 if _,err:=globalOnly.executionAdapter("dev",tpl,a,actor,call);!errors.Is(err,ErrExecutionUnavailable){t.Fatal("global grant accepted school request",err)}
}
