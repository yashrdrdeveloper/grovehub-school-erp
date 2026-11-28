// Package auditdelivery relays transactionally captured workflow events.
// The local database owns retries; the Audit provider owns deduplication.
package auditdelivery

import (
 "context"
 "database/sql"
 "errors"
 "log/slog"
 "time"

 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/client/audit"
)

type Worker struct{ db *sql.DB; publish audit.Publisher }
func New(db *sql.DB,publish audit.Publisher)*Worker{return &Worker{db:db,publish:publish}}

func(w *Worker) Step(ctx context.Context)(bool,error){
 ctx,cancel:=context.WithTimeout(ctx,10*time.Second);defer cancel()
 tx,err:=w.db.BeginTx(ctx,nil);if err!=nil{return false,err};defer tx.Rollback()
 var id int64;var a audit.Activity
 err=tx.QueryRowContext(ctx,"SELECT o.workflow_event_id,o.event_id,e.actor_id,e.resource_id,e.action,e.occurred_at,e.correlation_id FROM internal_tool.audit_outbox o JOIN internal_tool.workflow_events e ON e.id=o.workflow_event_id WHERE o.delivered_at IS NULL AND o.blocked_at IS NULL AND o.next_attempt_at<=now() ORDER BY o.next_attempt_at,o.workflow_event_id LIMIT 1 FOR UPDATE OF o SKIP LOCKED").Scan(&id,&a.EventID,&a.ActorID,&a.ResourceID,&a.Operation,&a.OccurredAt,&a.CorrelationID)
 if errors.Is(err,sql.ErrNoRows){return false,nil};if err!=nil{return false,err}
 call,stop:=context.WithTimeout(ctx,5*time.Second)
 status,sendErr:=w.publish(call,a);stop()
 switch{
 case sendErr==nil&&status==204:
  _,err=tx.ExecContext(ctx,"UPDATE internal_tool.audit_outbox SET delivered_at=now(),attempts=attempts+1 WHERE workflow_event_id=$1",id)
 case sendErr==nil&&(status==400||status==409):
  _,err=tx.ExecContext(ctx,"UPDATE internal_tool.audit_outbox SET blocked_at=now(),attempts=attempts+1 WHERE workflow_event_id=$1",id)
  slog.Warn("audit delivery blocked","workflow_event_id",id,"status",status)
 default:
  _,err=tx.ExecContext(ctx,"UPDATE internal_tool.audit_outbox SET attempts=attempts+1,next_attempt_at=now()+interval '30 seconds' WHERE workflow_event_id=$1",id)
 }
 if err!=nil{return false,err}
 return true,tx.Commit()
}
func(w *Worker) Run(ctx context.Context){
 for ctx.Err()==nil{
  worked,err:=w.Step(ctx)
  if err!=nil&&ctx.Err()==nil{slog.Warn("audit delivery step failed")}
  if worked&&err==nil{continue}
  select{case <-ctx.Done():return;case <-time.After(time.Second):}
 }
}
