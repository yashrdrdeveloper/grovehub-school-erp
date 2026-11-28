package workflow

import (
 "net/http"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
)
func(c *Controller) registerSearch(mux *http.ServeMux){
 bind:=func(path string,h handler){mux.HandleFunc("POST /internal/protected/"+path+"/search",c.auth.Protected(true,h))}
 bind("environments",func(w http.ResponseWriter,r *http.Request,a entity.Operator){in,ok:=decode[entity.PageRequest](w,r);if !ok{return};v,err:=c.svc.SearchEnvironments(r.Context(),in);respond(w,200,v,err)})
 bind("api-templates",func(w http.ResponseWriter,r *http.Request,a entity.Operator){in,ok:=decode[entity.PageRequest](w,r);if !ok{return};v,err:=c.svc.SearchTemplates(r.Context(),in);respond(w,200,v,err)})
 bind("requests",func(w http.ResponseWriter,r *http.Request,a entity.Operator){in,ok:=decode[entity.PageRequest](w,r);if !ok{return};v,err:=c.svc.SearchApprovals(r.Context(),a,in,false);respond(w,200,v,err)})
 bind("approvals",func(w http.ResponseWriter,r *http.Request,a entity.Operator){in,ok:=decode[entity.PageRequest](w,r);if !ok{return};v,err:=c.svc.SearchApprovals(r.Context(),a,in,true);respond(w,200,v,err)})
 bind("workflow-events",func(w http.ResponseWriter,r *http.Request,a entity.Operator){in,ok:=decode[entity.PageRequest](w,r);if !ok{return};v,err:=c.svc.SearchEvents(r.Context(),a,in);respond(w,200,v,err)})
}
