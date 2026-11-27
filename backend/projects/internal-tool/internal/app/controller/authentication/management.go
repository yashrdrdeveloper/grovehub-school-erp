package authentication

import (
 "database/sql"
 "encoding/json"
 "errors"
 "io"
 "net/http"

 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
 operators "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/repository/operator"
 service "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/services/authentication"
)

func managementResponse(w http.ResponseWriter,status int,data any,err error){
 if err==nil{write(w,status,map[string]any{"success":true,"data":data});return}
 code:="UNAVAILABLE";status=503
 switch{
 case errors.Is(err,operators.ErrForbidden):status=403;code="FORBIDDEN"
 case errors.Is(err,operators.ErrConflict):status=409;code="CONFLICT"
 case errors.Is(err,sql.ErrNoRows):status=404;code="NOT_FOUND"
 case errors.Is(err,service.ErrInvalid):status=400;code="INVALID_REQUEST"
 case errors.Is(err,service.ErrLimited):status=429;code="RATE_LIMITED"
 }
 write(w,status,map[string]any{"success":false,"errorCode":code})
}
func(c *Controller) registerManagement(mux *http.ServeMux){
 mux.HandleFunc("GET /internal/protected/operators",c.Protected(false,func(w http.ResponseWriter,r *http.Request,a entity.Operator){
  rows,err:=c.service.ListOperators(r.Context(),a);managementResponse(w,200,rows,err)
 }))
 mux.HandleFunc("POST /internal/protected/operators",c.Protected(true,func(w http.ResponseWriter,r *http.Request,a entity.Operator){
  var in struct{Email string `json:"email"`;Password string `json:"password"`;Persona entity.Persona `json:"persona"`}
  var extra any;d:=json.NewDecoder(http.MaxBytesReader(w,r.Body,4096));d.DisallowUnknownFields()
  if d.Decode(&in)!=nil||d.Decode(&extra)!=io.EOF{managementResponse(w,400,nil,service.ErrInvalid);return}
  row,err:=c.service.CreateOperator(r.Context(),a,in.Email,in.Password,in.Persona);managementResponse(w,201,row,err)
 }))
 for _,action:=range []struct{path,status string}{{"disable","DISABLED"},{"activate","ACTIVE"}}{
  mux.HandleFunc("POST /internal/protected/operators/{id}/"+action.path,c.Protected(true,func(w http.ResponseWriter,r *http.Request,a entity.Operator){
   raw,err:=io.ReadAll(http.MaxBytesReader(w,r.Body,1));if err!=nil||len(raw)!=0{managementResponse(w,400,nil,service.ErrInvalid);return}
   err=c.service.SetOperatorStatus(r.Context(),a,r.PathValue("id"),action.status);managementResponse(w,200,map[string]string{"status":action.status},err)
  }))
 }
}
