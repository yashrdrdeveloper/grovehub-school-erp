package auditdelivery

import("testing";"time")

func TestDeliveryHealth(t *testing.T){
 now:=time.Now()
 old:=now.Add(-BacklogThreshold)
 recent:=now.Add(-time.Second)
 for _,tc:=range []struct{s Summary;want string}{
  {Summary{},"DISABLED"},
  {Summary{Enabled:true},"HEALTHY"},
  {Summary{Enabled:true,Blocked:1},"DEGRADED"},
  {Summary{Enabled:true,Pending:1,OldestPendingAt:&old},"DEGRADED"},
  {Summary{Enabled:true,Pending:1,OldestPendingAt:&recent},"HEALTHY"},
  {Summary{Enabled:true,Pending:1},"DEGRADED"},
 }{if got:=DeliveryHealth(tc.s,now);got!=tc.want{t.Errorf("health=%s want %s",got,tc.want)}}
}
