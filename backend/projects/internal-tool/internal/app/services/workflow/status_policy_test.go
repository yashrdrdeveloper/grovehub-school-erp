package workflow

import (
 "context"
 "errors"
 "testing"
 "github.com/google/uuid"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
)
func TestDisabledSchoolLifecycleGrant(t *testing.T){
 school,revision,otherRevision,requester,reviewer:=uuid.NewString(),uuid.NewString(),uuid.NewString(),uuid.NewString(),uuid.NewString()
 base:=SchoolBinding{Mutation:true,Operators:[]string{requester,reviewer},Validate:func(Invocation)bool{return true},Adapter:func(context.Context,Invocation)(int,error){return 204,nil}}
 lifecycle:=base;lifecycle.Lifecycle=true
 s:=NewWithSchoolBindings(nil,nil,map[SchoolBindingKey]SchoolBinding{{"dev",revision,school}:lifecycle,{"dev",otherRevision,school}:base})
 s.WithDirectory(func(context.Context,[]string)(map[string]SchoolRecord,error){return map[string]SchoolRecord{school:{ID:school,Name:"School",Status:"DISABLED",Version:2}},nil})
 ctx:=context.Background();actor:=entity.Operator{ID:requester,Persona:entity.L2}
 grants,err:=s.ResolvedSchoolGrants(ctx,actor)
 if err!=nil||len(grants)!=1||grants[0].TemplateID!=revision||grants[0].SchoolStatus!="DISABLED"||grants[0].SchoolVersion!="2"{t.Fatal("incorrect disabled grants",grants,err)}
 if err=s.checkDirectoryForBinding(ctx,school,base);!errors.Is(err,ErrForbidden){t.Fatal("disabled rename allowed",err)}
 if err=s.checkDirectoryForBinding(ctx,school,lifecycle);err!=nil{t.Fatal("reactivation denied",err)}
 tplt:=entity.Template{ID:revision,Scope:"SCHOOL",Method:"PATCH"}
 a:=entity.Approval{SchoolID:&school,RequesterID:requester,ReviewerID:&reviewer}
 call:=Invocation{SchoolID:&school,Method:"PATCH"}
 if _,err=s.executionAdapter("dev",tplt,a,actor,call);err!=nil{t.Fatal("status binding denied",err)}
 a.ReviewerID=&requester
 if _,err=s.executionAdapter("dev",tplt,a,actor,call);!errors.Is(err,ErrForbidden){t.Fatal("self approval allowed",err)}
}
