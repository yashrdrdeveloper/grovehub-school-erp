package invitation

import (
 "errors"
 "testing"
 schoolgateway "grovehub/libs/platform/session"
)

func TestTenantDeliveryDecision(t *testing.T) {
 for _,tc := range []struct{
  name string
  state schoolgateway.TenantState
  issued int64
  err error
  reason string
  retry bool
 }{
  {"active current",schoolgateway.TenantState{Generation:2,Active:true},2,nil,"",false},
  {"reactivated stale",schoolgateway.TenantState{Generation:2,Active:true},1,nil,"INVALID_PROOF",false},
  {"legacy",schoolgateway.TenantState{Generation:2,Active:true},0,nil,"INVALID_PROOF",false},
  {"disabled",schoolgateway.TenantState{Generation:2},2,nil,"ACCOUNT_INELIGIBLE",false},
  {"unavailable",schoolgateway.TenantState{},2,errors.New("network unavailable"),"",true},
  {"invalid state",schoolgateway.TenantState{Active:true},2,nil,"",true},
 }{
  t.Run(tc.name,func(t *testing.T){
   reason,err:=tenantDeliveryDecision(tc.state,tc.err,tc.issued)
   if reason!=tc.reason||(err!=nil)!=tc.retry{t.Fatal(reason,err)}
  })
 }
}
