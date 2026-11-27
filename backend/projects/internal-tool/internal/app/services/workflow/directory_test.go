package workflow

import (
 "context"
 "errors"
 "testing"

 "github.com/google/uuid"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
)
func TestResolvedSchoolGrants(t *testing.T){
 actor,active,disabled,missing,other:=uuid.NewString(),uuid.NewString(),uuid.NewString(),uuid.NewString(),uuid.NewString()
 bindings:=map[SchoolBindingKey]SchoolBinding{}
 for _,id:=range []string{active,disabled,missing,other}{
  operator:=actor;if id==other{operator=other}
  bindings[SchoolBindingKey{"dev",uuid.NewString(),id}]=SchoolBinding{Operators:[]string{operator},Validate:func(Invocation)bool{return true},Adapter:func(context.Context,Invocation)(int,error){return 204,nil}}
 }
 svc:=NewWithSchoolBindings(nil,nil,bindings).WithDirectory(func(_ context.Context,ids []string)(map[string]SchoolRecord,error){
  for _,id:=range ids{if id==other{t.Fatal("unassigned school sent to directory")}}
  return map[string]SchoolRecord{active:{ID:active,Name:"Active school",Code:"active",Status:"ACTIVE"},disabled:{ID:disabled,Status:"DISABLED"}},nil
 })
 rows,err:=svc.ResolvedSchoolGrants(context.Background(),entity.Operator{ID:actor,Persona:entity.L1})
 if err!=nil||len(rows)!=1||rows[0].SchoolID!=active||rows[0].SchoolName!="Active school"{t.Fatal("directory resolution failed",rows,err)}
 if err=svc.checkDirectory(context.Background(),disabled);!errors.Is(err,ErrForbidden){t.Fatal("disabled school permitted",err)}
 if err=svc.checkDirectory(context.Background(),missing);!errors.Is(err,ErrForbidden){t.Fatal("missing school permitted",err)}
 svc.WithDirectory(func(context.Context,[]string)(map[string]SchoolRecord,error){return nil,errors.New("directory offline")})
 if _,err=svc.ResolvedSchoolGrants(context.Background(),entity.Operator{ID:actor,Persona:entity.L1});err==nil{t.Fatal("directory outage bypassed")}
}
