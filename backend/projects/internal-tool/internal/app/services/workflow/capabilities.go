package workflow

import (
 "context"
 "database/sql"
 "errors"

 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
)

// RequestCapabilities describes actions available for the current request
// snapshot. Mutating routes always repeat authorization and version checks.
type RequestCapabilities struct {
 Version int64 `json:"version"`
 CanApprove bool `json:"canApprove"`
 CanReject bool `json:"canReject"`
 CanCancel bool `json:"canCancel"`
 CanExecute bool `json:"canExecute"`
 CanReconcile bool `json:"canReconcile"`
}

func(s *Service) Capabilities(ctx context.Context,actor entity.Operator,id string)(RequestCapabilities,error){
 a,err:=s.Get(ctx,actor,id)
 if err!=nil{return RequestCapabilities{},err}
 result:=RequestCapabilities{Version:a.Version}
 if a.Status=="PENDING"{
  result.CanCancel=actor.ID==a.RequesterID
  result.CanReject=actor.Persona==entity.L2
  if result.CanReject{
   err=s.checkSchoolReview(ctx,actor,id)
   if err!=nil&&!errors.Is(err,ErrForbidden){return RequestCapabilities{},err}
   result.CanApprove=err==nil
   if result.CanApprove&&a.SchoolID!=nil&&a.RequesterID==actor.ID{
    t,err:=s.repo.Template(ctx,a.TemplateID);if err!=nil{return RequestCapabilities{},err}
    result.CanApprove=t.Method=="GET"
   }
  }
  return result,nil
 }
 if a.Status!="APPROVED"{return result,nil}
 e,err:=s.repo.Execution(ctx,actor,id)
 if err==nil{
  if actor.Persona==entity.L2&&(e.Status=="RUNNING"||e.Status=="UNKNOWN"){
   _,err=s.Reconciliation(ctx,actor,id)
   if err!=nil&&!errors.Is(err,sql.ErrNoRows){return RequestCapabilities{},err}
   // The actual reconciliation transaction enforces the stale RUNNING cutoff.
   result.CanReconcile=errors.Is(err,sql.ErrNoRows)
  }
  return result,nil
 }
 if !errors.Is(err,sql.ErrNoRows){return RequestCapabilities{},err}
 t,err:=s.repo.Template(ctx,a.TemplateID);if err!=nil{return RequestCapabilities{},err}
 env,err:=s.repo.Environment(ctx,a.EnvironmentID);if err!=nil{return RequestCapabilities{},err}
 invocation,err:=resolve(t,a);if err!=nil{return RequestCapabilities{},err}
 _,err=s.executionAdapter(env.Code,t,a,actor,invocation)
 if err!=nil&&!errors.Is(err,ErrForbidden)&&!errors.Is(err,ErrExecutionUnavailable){return RequestCapabilities{},err}
 result.CanExecute=err==nil
 if result.CanExecute&&a.SchoolID!=nil{
  if err=s.checkDirectory(ctx,*a.SchoolID);err!=nil{
   if !errors.Is(err,ErrForbidden){return RequestCapabilities{},err}
   result.CanExecute=false
  }
 }
 return result,nil
}
