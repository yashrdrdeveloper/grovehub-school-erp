package delivery

import (
 "context"
 "database/sql"
 "net/http"
 "time"
)

// Health includes unconfigured schools: committed receipts must not silently
// accumulate merely because their publisher has not been provisioned.
func Health(db *sql.DB,configs []Config)http.HandlerFunc{
 return func(w http.ResponseWriter,r *http.Request){
  w.Header().Set("Cache-Control","no-store")
  if len(configs)==0{w.WriteHeader(503);return}
  ctx,cancel:=context.WithTimeout(r.Context(),2*time.Second);defer cancel()
  ids:=make([]string,0,len(configs));for _,c:=range configs{ids=append(ids,c.SchoolID)}
  var degraded bool
  err:=db.QueryRowContext(ctx,`SELECT EXISTS(
   SELECT 1 FROM school_registry.rename_outbox o JOIN school_registry.rename_events e USING(execution_id)
   WHERE o.delivered_at IS NULL AND
    (o.blocked_at IS NOT NULL OR e.occurred_at<=now()-interval '5 minutes' OR NOT(e.school_id=ANY($1::uuid[])))
  )`,ids).Scan(&degraded)
  if err!=nil||degraded{w.WriteHeader(503);return}
  w.WriteHeader(204)
 }
}
