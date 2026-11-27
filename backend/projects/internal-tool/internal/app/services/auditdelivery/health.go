package auditdelivery

import (
 "context"
 "net/http"
 "time"
)

const BacklogThreshold=5*time.Minute

// DeliveryHealth is deliberately separate from application readiness.
// A disabled optional relay is distinguishable from an unhealthy relay.
func DeliveryHealth(s Summary,now time.Time)string{
 if !s.Enabled{return "DISABLED"}
 if s.Blocked>0{return "DEGRADED"}
 if s.Pending>0&&(s.OldestPendingAt==nil||now.Sub(*s.OldestPendingAt)>=BacklogThreshold){return "DEGRADED"}
 return "HEALTHY"
}
func(m *Monitor) HealthHandler(w http.ResponseWriter,r *http.Request){
 w.Header().Set("Cache-Control","no-store")
 ctx,cancel:=context.WithTimeout(r.Context(),2*time.Second);defer cancel()
 summary,err:=m.Summary(ctx)
 if err!=nil||summary.Health!="HEALTHY"{w.WriteHeader(http.StatusServiceUnavailable);return}
 w.WriteHeader(http.StatusNoContent)
}
