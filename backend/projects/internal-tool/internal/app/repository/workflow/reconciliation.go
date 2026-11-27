package workflow

import (
 "context"
 "database/sql"

 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
)

func(r *Repository) Reconcile(ctx context.Context,actor entity.Operator,requestID,outcome,evidence string)(entity.Reconciliation,error){
 if actor.Persona!=entity.L2{return entity.Reconciliation{},ErrForbidden}
 tx,err:=r.db.BeginTx(ctx,nil);if err!=nil{return entity.Reconciliation{},err};defer tx.Rollback()
 var active bool
 if err=tx.QueryRowContext(ctx,"SELECT EXISTS(SELECT 1 FROM internal_tool.operators WHERE id=$1 AND persona='L2' AND status='ACTIVE')",actor.ID).Scan(&active);err!=nil{return entity.Reconciliation{},err}
 if !active{return entity.Reconciliation{},ErrForbidden}
 var id,status string;var stale bool
 if err=tx.QueryRowContext(ctx,"SELECT id,status,created_at<now()-interval '5 minutes' FROM internal_tool.executions WHERE request_id=$1 FOR UPDATE",requestID).Scan(&id,&status,&stale);err!=nil{return entity.Reconciliation{},err}
 if status=="RUNNING"&&stale{
  if _,err=tx.ExecContext(ctx,"UPDATE internal_tool.executions SET status='UNKNOWN',finished_at=now() WHERE id=$1",id);err!=nil{return entity.Reconciliation{},err}
  if err=event(ctx,tx,actor.ID,requestID,"EXECUTION_ABANDONED");err!=nil{return entity.Reconciliation{},err}
  status="UNKNOWN"
 }
 if status!="UNKNOWN"{return entity.Reconciliation{},ErrConflict}
 var row entity.Reconciliation
 err=tx.QueryRowContext(ctx,"INSERT INTO internal_tool.execution_reconciliations(execution_id,reviewer_id,outcome,evidence) VALUES($1,$2,$3,$4) ON CONFLICT(execution_id) DO NOTHING RETURNING execution_id,reviewer_id,outcome,evidence,created_at",id,actor.ID,outcome,evidence).Scan(&row.ExecutionID,&row.ReviewerID,&row.Outcome,&row.Evidence,&row.CreatedAt)
 if err==sql.ErrNoRows{return row,ErrConflict};if err!=nil{return row,err}
 if err=event(ctx,tx,actor.ID,requestID,"EXECUTION_RECONCILED");err!=nil{return row,err}
 return row,tx.Commit()
}

func(r *Repository) Reconciliation(ctx context.Context,actor entity.Operator,requestID string)(entity.Reconciliation,error){
 if _,err:=r.Get(ctx,actor,requestID);err!=nil{return entity.Reconciliation{},err}
 var row entity.Reconciliation
 err:=r.db.QueryRowContext(ctx,"SELECT r.execution_id,r.reviewer_id,r.outcome,r.evidence,r.created_at FROM internal_tool.execution_reconciliations r JOIN internal_tool.executions e ON e.id=r.execution_id WHERE e.request_id=$1",requestID).Scan(&row.ExecutionID,&row.ReviewerID,&row.Outcome,&row.Evidence,&row.CreatedAt)
 return row,err
}
