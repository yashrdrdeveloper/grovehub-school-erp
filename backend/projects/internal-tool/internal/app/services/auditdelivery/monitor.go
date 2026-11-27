package auditdelivery

import (
 "context"
 "database/sql"
 "errors"
 "time"

 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
)
var ErrConflict=errors.New("audit delivery not blocked")
type Summary struct{
 Enabled bool `json:"enabled"`
 Pending int64 `json:"pending"`
 Blocked int64 `json:"blocked"`
 Delivered int64 `json:"delivered"`
 OldestPendingAt *time.Time `json:"oldestPendingAt"`
 LastDeliveredAt *time.Time `json:"lastDeliveredAt"`
}
type BlockedEvent struct{
 EventID string `json:"eventId"`
 ActorID string `json:"actorId"`
 ResourceID string `json:"resourceId"`
 Action string `json:"action"`
 Attempts int `json:"attempts"`
 BlockedAt time.Time `json:"blockedAt"`
}
type Monitor struct{db *sql.DB;enabled bool}
func NewMonitor(db *sql.DB,enabled bool)*Monitor{return &Monitor{db:db,enabled:enabled}}
func(m *Monitor) Summary(ctx context.Context)(Summary,error){
 v:=Summary{Enabled:m.enabled}
 err:=m.db.QueryRowContext(ctx,"SELECT count(*) FILTER(WHERE delivered_at IS NULL AND blocked_at IS NULL), count(*) FILTER(WHERE blocked_at IS NOT NULL), count(*) FILTER(WHERE delivered_at IS NOT NULL), min(e.occurred_at) FILTER(WHERE o.delivered_at IS NULL AND o.blocked_at IS NULL), max(o.delivered_at) FROM internal_tool.audit_outbox o JOIN internal_tool.workflow_events e ON e.id=o.workflow_event_id").Scan(&v.Pending,&v.Blocked,&v.Delivered,&v.OldestPendingAt,&v.LastDeliveredAt)
 return v,err
}
func(m *Monitor) Blocked(ctx context.Context,p entity.PageRequest)(entity.Page[BlockedEvent],error){
 rows,err:=m.db.QueryContext(ctx,"SELECT o.event_id,e.actor_id,e.resource_id,e.action,o.attempts,o.blocked_at FROM internal_tool.audit_outbox o JOIN internal_tool.workflow_events e ON e.id=o.workflow_event_id WHERE o.blocked_at IS NOT NULL ORDER BY o.blocked_at DESC,o.event_id DESC LIMIT $1 OFFSET $2",p.RowSize+1,p.Offset())
 if err!=nil{return entity.Page[BlockedEvent]{},err};defer rows.Close()
 records:=[]BlockedEvent{}
 for rows.Next(){var v BlockedEvent;if err=rows.Scan(&v.EventID,&v.ActorID,&v.ResourceID,&v.Action,&v.Attempts,&v.BlockedAt);err!=nil{return entity.Page[BlockedEvent]{},err};records=append(records,v)}
 if err=rows.Err();err!=nil{return entity.Page[BlockedEvent]{},err}
 return entity.NewPage(p,records),nil
}
// Retry changes delivery state only, not the immutable occurrence or its UUID.
// The requeue action is itself captured by the workflow outbox trigger.
func(m *Monitor) Retry(ctx context.Context,actor,id string)error{
 tx,err:=m.db.BeginTx(ctx,nil);if err!=nil{return err};defer tx.Rollback()
 var active bool
 if err=tx.QueryRowContext(ctx,"SELECT EXISTS(SELECT 1 FROM internal_tool.operators WHERE id=$1 AND persona='L2' AND status='ACTIVE')",actor).Scan(&active);err!=nil{return err}
 if !active{return ErrConflict}
 result,err:=tx.ExecContext(ctx,"UPDATE internal_tool.audit_outbox SET blocked_at=NULL,next_attempt_at=now() WHERE event_id=$1 AND blocked_at IS NOT NULL AND delivered_at IS NULL",id)
 if err!=nil{return err};n,err:=result.RowsAffected();if err!=nil{return err};if n!=1{return ErrConflict}
 if _,err=tx.ExecContext(ctx,"INSERT INTO internal_tool.workflow_events(actor_id,resource_id,action) VALUES($1,$2,'AUDIT_DELIVERY_REQUEUED')",actor,id);err!=nil{return err}
 return tx.Commit()
}
