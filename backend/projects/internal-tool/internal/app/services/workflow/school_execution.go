package workflow

import (
 "database/sql"

 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
)

// SchoolBindingKey binds an immutable template revision to an exact school and
// logical environment. Global bindings cannot dispatch school requests.
type SchoolBindingKey struct{EnvironmentCode,TemplateID,SchoolID string}
type SchoolBinding struct{
 Adapter Adapter
 Validate func(Invocation)bool
 Operators []string
}
func NewWithSchoolBindings(db *sql.DB,global map[BindingKey]Adapter,school map[SchoolBindingKey]SchoolBinding)*Service{
 s:=NewWithBindings(db,global)
 s.schoolBindings=make(map[SchoolBindingKey]SchoolBinding,len(school))
 for key,binding:=range school{
  if binding.Adapter==nil||binding.Validate==nil||len(binding.Operators)==0||!validID(key.SchoolID)||!validID(key.TemplateID)||!codePattern.MatchString(key.EnvironmentCode){continue}
  binding.Operators=append([]string(nil),binding.Operators...)
  s.schoolBindings[key]=binding
 }
 return s
}
func(s *Service) executionAdapter(environment string,t entity.Template,a entity.Approval,actor entity.Operator,call Invocation)(Adapter,error){
 if t.Scope=="GLOBAL"&&a.SchoolID==nil{
  adapter:=s.bindings[BindingKey{environment,t.ID}]
  if adapter==nil{return nil,ErrExecutionUnavailable}
  return adapter,nil
 }
 if t.Scope!="SCHOOL"||a.SchoolID==nil||!validID(*a.SchoolID){return nil,ErrForbidden}
 // Until a separately reviewed mutation policy is implemented, school
 // execution is read-only, regardless of what a template declares.
 if t.Method!="GET"||call.Method!="GET"{return nil,ErrForbidden}
 binding,ok:=s.schoolBindings[SchoolBindingKey{environment,t.ID,*a.SchoolID}]
 if !ok{return nil,ErrExecutionUnavailable}
 allowed:=func(id string)bool{
  for _,operator:=range binding.Operators{if operator==id{return true}}
  return false
 }
 if !allowed(actor.ID)||!allowed(a.RequesterID)||a.ReviewerID==nil||!allowed(*a.ReviewerID)||!binding.Validate(call){return nil,ErrForbidden}
 return binding.Adapter,nil
}
