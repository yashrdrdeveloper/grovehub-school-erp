package audit

import (
 "context"
 "errors"
 "net/http"
 core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
 sdk "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/generated/audit"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
)
type Reader func(context.Context,entity.PageRequest)(*sdk.OperatorActivityPage,error)
func NewReader(c Config)(Reader,error){
 api,err:=buildClient(c,"internal-tool-audit-reader",readPolicy)
 if err!=nil||api==nil{return nil,err}
 return func(ctx context.Context,p entity.PageRequest)(*sdk.OperatorActivityPage,error){
  page,response,err:=api.DefaultAPI.SearchOperatorActivities(ctx).OperatorSearchInput(sdk.OperatorSearchInput{PageNumber:int32(p.PageNumber),RowSize:int32(p.RowSize)}).Execute()
  if response!=nil&&response.Body!=nil{defer response.Body.Close()}
  if err!=nil||response==nil||response.StatusCode!=200||page==nil{return nil,errors.New("central audit unavailable")}
  return page,nil
 },nil
}
func readPolicy(r *http.Request)(core.WorkloadClaims,error){
 if r==nil||r.URL==nil||r.Method!="POST"||r.URL.Path!="/audit/restricted/operator-activities/search"||r.URL.RawPath!=""||r.URL.RawQuery!=""||r.URL.ForceQuery||r.Header.Get("Cookie")!=""||r.Header.Get("X-GH-Identity-Proof")!=""{return core.WorkloadClaims{},core.ErrWorkloadDenied}
 return core.WorkloadClaims{Scopes:[]string{"internal.audit.operator.read"}},nil
}
