package workflow

import (
 "database/sql"
 "sort"

 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
)

// SchoolBindingKey binds an immutable template revision to an exact school and
// logical environment. Global bindings cannot dispatch school requests.
type SchoolBindingKey struct{EnvironmentCode,TemplateID,SchoolID string}
type SchoolBinding struct{
 Adapter Adapter
 Validate func(Invocation)bool
 Operators []string
 Requesters,Reviewers,Executors []string
}
func NewWithSchoolBindings(db *sql.DB,global map[BindingKey]Adapter,school map[SchoolBindingKey]SchoolBinding)*Service{
 s:=NewWithBindings(db,global)
 s.schoolBindings=make(map[SchoolBindingKey]SchoolBinding,len(school))
 for key,binding:=range school{
  if binding.Adapter==nil||binding.Validate==nil||len(binding.Operators)==0||!validID(key.SchoolID)||!validID(key.TemplateID)||!codePattern.MatchString(key.EnvironmentCode){continue}
  binding.Operators=append([]string(nil),binding.Operators...)
  if len(binding.Requesters)==0&&len(binding.Reviewers)==0&&len(binding.Executors)==0{binding.Requesters=binding.Operators;binding.Reviewers=binding.Operators;binding.Executors=binding.Operators}
  if len(binding.Requesters)==0||len(binding.Reviewers)==0||len(binding.Executors)==0{continue}
  binding.Requesters=append([]string(nil),binding.Requesters...);binding.Reviewers=append([]string(nil),binding.Reviewers...);binding.Executors=append([]string(nil),binding.Executors...)
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
 if !allowedOperator(binding.Executors,actor.ID)||!allowedOperator(binding.Requesters,a.RequesterID)||a.ReviewerID==nil||!allowedOperator(binding.Reviewers,*a.ReviewerID)||!binding.Validate(call){return nil,ErrForbidden}
 return binding.Adapter,nil
}

func allowedOperator(ids []string,id string)bool{
 for _,candidate:=range ids{if candidate==id{return true}}
 return false
}

type SchoolGrantView struct{
 EnvironmentCode string `json:"environmentCode"`
 TemplateID string `json:"templateId"`
 SchoolID string `json:"schoolId"`
 CanRequest bool `json:"canRequest"`
 CanReview bool `json:"canReview"`
 CanExecute bool `json:"canExecute"`
}

// SchoolGrants lists only the authenticated operator's configured capabilities.
// No credentials, origin, or other operator identities are exposed.
func(s *Service) SchoolGrants(actor entity.Operator)[]SchoolGrantView{
 result:=[]SchoolGrantView{}
 for key,b:=range s.schoolBindings{
  v:=SchoolGrantView{EnvironmentCode:key.EnvironmentCode,TemplateID:key.TemplateID,SchoolID:key.SchoolID,CanRequest:allowedOperator(b.Requesters,actor.ID),CanReview:actor.Persona==entity.L2&&allowedOperator(b.Reviewers,actor.ID),CanExecute:allowedOperator(b.Executors,actor.ID)}
  if v.CanRequest||v.CanReview||v.CanExecute{result=append(result,v)}
 }
 sort.Slice(result,func(i,j int)bool{
  a,b:=result[i],result[j]
  if a.EnvironmentCode!=b.EnvironmentCode{return a.EnvironmentCode<b.EnvironmentCode}
  if a.SchoolID!=b.SchoolID{return a.SchoolID<b.SchoolID}
  return a.TemplateID<b.TemplateID
 })
 return result
}
