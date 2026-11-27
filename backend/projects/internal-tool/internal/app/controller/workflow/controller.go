package workflow

import (
 "database/sql"
 "encoding/json"
 "errors"
 "io"
 "net/http"

 auth "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/controller/authentication"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
 service "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/services/workflow"
)

type Controller struct{svc *service.Service;auth *auth.Controller}
func New(s *service.Service,a *auth.Controller)*Controller{return &Controller{s,a}}
func respond(w http.ResponseWriter,status int,data any,err error){
 w.Header().Set("Content-Type","application/json");w.Header().Set("Cache-Control","no-store")
 if err!=nil{
  code:="UNAVAILABLE";status=503
  switch{case errors.Is(err,service.ErrInvalid):status=400;code="INVALID_REQUEST";case errors.Is(err,service.ErrForbidden):status=403;code="FORBIDDEN";case errors.Is(err,service.ErrConflict):status=409;code="CONFLICT";case errors.Is(err,sql.ErrNoRows):status=404;code="NOT_FOUND";case errors.Is(err,service.ErrExecutionUnavailable):status=503;code="EXECUTION_UNAVAILABLE"}
  w.WriteHeader(status);_ =json.NewEncoder(w).Encode(map[string]any{"success":false,"errorCode":code});return
 }
 w.WriteHeader(status);_ =json.NewEncoder(w).Encode(map[string]any{"success":true,"data":data})
}
func decode[T any](w http.ResponseWriter,r *http.Request)(T,bool){
 var in T;var extra any
 d:=json.NewDecoder(http.MaxBytesReader(w,r.Body,32768));d.DisallowUnknownFields()
 if d.Decode(&in)!=nil||d.Decode(&extra)!=io.EOF{respond(w,400,nil,service.ErrInvalid);return in,false};return in,true
}
type handler func(http.ResponseWriter,*http.Request,entity.Operator)
func(c *Controller) Register(mux *http.ServeMux){
 c.registerSearch(mux)
 bind:=func(pattern string,mutate bool,h handler){mux.HandleFunc(pattern,c.auth.Protected(mutate,h))}
 bind("POST /internal/protected/environments",true,func(w http.ResponseWriter,r *http.Request,a entity.Operator){in,ok:=decode[service.EnvironmentInput](w,r);if !ok{return};v,err:=c.svc.CreateEnvironment(r.Context(),a,in);respond(w,201,v,err)})
 bind("GET /internal/protected/school-grants",false,func(w http.ResponseWriter,r *http.Request,a entity.Operator){respond(w,200,c.svc.SchoolGrants(a),nil)})
 bind("GET /internal/protected/environments",false,func(w http.ResponseWriter,r *http.Request,a entity.Operator){v,err:=c.svc.Environments(r.Context());respond(w,200,v,err)})
 bind("POST /internal/protected/api-templates",true,func(w http.ResponseWriter,r *http.Request,a entity.Operator){in,ok:=decode[service.TemplateInput](w,r);if !ok{return};v,err:=c.svc.CreateTemplate(r.Context(),a,in);respond(w,201,v,err)})
 bind("GET /internal/protected/api-templates",false,func(w http.ResponseWriter,r *http.Request,a entity.Operator){v,err:=c.svc.Templates(r.Context());respond(w,200,v,err)})
 bind("POST /internal/protected/requests",true,func(w http.ResponseWriter,r *http.Request,a entity.Operator){in,ok:=decode[service.SubmitInput](w,r);if !ok{return};v,err:=c.svc.Submit(r.Context(),a,in);respond(w,201,v,err)})
 bind("GET /internal/protected/requests",false,func(w http.ResponseWriter,r *http.Request,a entity.Operator){v,err:=c.svc.List(r.Context(),a,false);respond(w,200,v,err)})
 bind("POST /internal/protected/approvals/bulk-approve",true,func(w http.ResponseWriter,r *http.Request,a entity.Operator){in,ok:=decode[service.BulkApprovalInput](w,r);if !ok{return};v,err:=c.svc.BulkApprove(r.Context(),a,in);respond(w,200,v,err)})
 bind("GET /internal/protected/approvals",false,func(w http.ResponseWriter,r *http.Request,a entity.Operator){v,err:=c.svc.List(r.Context(),a,true);respond(w,200,v,err)})
 bind("GET /internal/protected/requests/{id}",false,func(w http.ResponseWriter,r *http.Request,a entity.Operator){v,err:=c.svc.Get(r.Context(),a,r.PathValue("id"));respond(w,200,v,err)})
 bind("POST /internal/protected/requests/{id}/execute",true,func(w http.ResponseWriter,r *http.Request,a entity.Operator){in,ok:=decode[service.ExecutionInput](w,r);if !ok{return};v,err:=c.svc.Execute(r.Context(),a,r.PathValue("id"),in);respond(w,200,v,err)})
 bind("GET /internal/protected/requests/{id}/execution",false,func(w http.ResponseWriter,r *http.Request,a entity.Operator){v,err:=c.svc.Execution(r.Context(),a,r.PathValue("id"));respond(w,200,v,err)})
 bind("POST /internal/protected/requests/{id}/reconcile",true,func(w http.ResponseWriter,r *http.Request,a entity.Operator){in,ok:=decode[service.ReconcileInput](w,r);if !ok{return};v,err:=c.svc.Reconcile(r.Context(),a,r.PathValue("id"),in);respond(w,200,v,err)})
 bind("GET /internal/protected/requests/{id}/reconciliation",false,func(w http.ResponseWriter,r *http.Request,a entity.Operator){v,err:=c.svc.Reconciliation(r.Context(),a,r.PathValue("id"));respond(w,200,v,err)})
 for _,action:=range []struct{path,status string}{{"approve","APPROVED"},{"reject","REJECTED"},{"cancel","CANCELLED"}}{
  bind("POST /internal/protected/requests/{id}/"+action.path,true,func(w http.ResponseWriter,r *http.Request,a entity.Operator){in,ok:=decode[service.DecisionInput](w,r);if !ok{return};v,err:=c.svc.Decide(r.Context(),a,r.PathValue("id"),action.status,in);respond(w,200,v,err)})
 }
}
