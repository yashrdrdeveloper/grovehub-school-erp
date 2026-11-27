package audit

import (
 "context"
 "errors"
 "net/http"
 core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
 sdk "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/generated/audit"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
)
type ReceiptReader func(context.Context,entity.PageRequest)(*sdk.SchoolReceiptPage,error)
func NewReceiptReader(c Config)(ReceiptReader,error){
 api,err:=buildClient(c,"internal-tool-school-receipt-reader",receiptReadPolicy)
 if err!=nil||api==nil{return nil,err}
 return func(ctx context.Context,p entity.PageRequest)(*sdk.SchoolReceiptPage,error){
  page,response,err:=api.DefaultAPI.SearchSchoolMutationReceipts(ctx).OperatorSearchInput(sdk.OperatorSearchInput{PageNumber:int32(p.PageNumber),RowSize:int32(p.RowSize)}).Execute()
  if response!=nil&&response.Body!=nil{defer response.Body.Close()}
  if err!=nil||response==nil||response.StatusCode!=200||page==nil{return nil,errors.New("central audit unavailable")}
  return page,nil
 },nil
}
func receiptReadPolicy(r *http.Request)(core.WorkloadClaims,error){
 if r==nil||r.URL==nil||r.Method!="POST"||r.URL.Path!="/audit/restricted/school-mutation-receipts/search"||r.URL.RawPath!=""||r.URL.RawQuery!=""||r.URL.ForceQuery||r.Header.Get("Cookie")!=""||r.Header.Get("X-GH-Identity-Proof")!=""{return core.WorkloadClaims{},core.ErrWorkloadDenied}
 return core.WorkloadClaims{Scopes:[]string{"internal.audit.school.mutation-receipt.read"}},nil
}
