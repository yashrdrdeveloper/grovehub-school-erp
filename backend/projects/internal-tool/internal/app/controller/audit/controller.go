package audit

import (
 "context"
 "encoding/json"
 "errors"
 "io"
 "net/http"
 "time"

 "github.com/google/uuid"
 auth "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/controller/authentication"
 client "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/client/audit"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/services/auditdelivery"
)
func response(w http.ResponseWriter,status int,data any){
 w.Header().Set("Content-Type","application/json");w.Header().Set("Cache-Control","no-store")
 w.WriteHeader(status)
 if status==200{_ =json.NewEncoder(w).Encode(map[string]any{"success":true,"data":data});return}
 code:="UNAVAILABLE"
 switch status{case 400:code="INVALID_REQUEST";case 403:code="FORBIDDEN";case 409:code="CONFLICT"}
 _ =json.NewEncoder(w).Encode(map[string]any{"success":false,"errorCode":code})
}
func page(w http.ResponseWriter,r *http.Request)(entity.PageRequest,bool){
 var p entity.PageRequest;var extra any
 d:=json.NewDecoder(http.MaxBytesReader(w,r.Body,4096));d.DisallowUnknownFields()
 if d.Decode(&p)!=nil||d.Decode(&extra)!=io.EOF{return p,false}
 return p.Normalize()
}
func Register(mux *http.ServeMux,a *auth.Controller,monitor *auditdelivery.Monitor,reader client.Reader,receiptReaders ...client.ReceiptReader){
 bind:=func(pattern string,mutate bool,h func(http.ResponseWriter,*http.Request,entity.Operator)){
  mux.HandleFunc(pattern,a.Protected(mutate,func(w http.ResponseWriter,r *http.Request,op entity.Operator){
   if op.Persona!=entity.L2{response(w,403,nil);return}
   ctx,cancel:=context.WithTimeout(r.Context(),5*time.Second);defer cancel()
   h(w,r.WithContext(ctx),op)
  }))
 }
 bind("POST /internal/protected/operator-activities/search",true,func(w http.ResponseWriter,r *http.Request,_ entity.Operator){
  p,ok:=page(w,r);if !ok{response(w,400,nil);return}
  if reader==nil{response(w,503,nil);return}
  v,err:=reader(r.Context(),p);if err!=nil{response(w,503,nil);return};response(w,200,v)
 })
 bind("POST /internal/protected/school-rename-receipts/search",true,func(w http.ResponseWriter,r *http.Request,_ entity.Operator){
  p,ok:=page(w,r);if !ok{response(w,400,nil);return}
  if len(receiptReaders)!=1||receiptReaders[0]==nil{response(w,503,nil);return}
  v,err:=receiptReaders[0](r.Context(),p);if err!=nil{response(w,503,nil);return};response(w,200,v)
 })
 bind("GET /internal/protected/audit-delivery/summary",false,func(w http.ResponseWriter,r *http.Request,_ entity.Operator){
  v,err:=monitor.Summary(r.Context());if err!=nil{response(w,503,nil);return};response(w,200,v)
 })
 bind("POST /internal/protected/audit-blocked/search",true,func(w http.ResponseWriter,r *http.Request,_ entity.Operator){
  p,ok:=page(w,r);if !ok{response(w,400,nil);return}
  v,err:=monitor.Blocked(r.Context(),p);if err!=nil{response(w,503,nil);return};response(w,200,v)
 })
 bind("POST /internal/protected/audit-blocked/{id}/retry",true,func(w http.ResponseWriter,r *http.Request,op entity.Operator){
  raw,err:=io.ReadAll(http.MaxBytesReader(w,r.Body,1));id:=r.PathValue("id");u,parseErr:=uuid.Parse(id)
  if err!=nil||len(raw)!=0||parseErr!=nil||u==uuid.Nil||u.String()!=id{response(w,400,nil);return}
  err=monitor.Retry(r.Context(),op.ID,id)
  if errors.Is(err,auditdelivery.ErrConflict){response(w,409,nil);return}
  if err!=nil{response(w,503,nil);return};response(w,200,map[string]string{"status":"REQUEUED"})
 })
}
