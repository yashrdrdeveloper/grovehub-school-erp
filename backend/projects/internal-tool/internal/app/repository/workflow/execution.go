package workflow

import (
 "context"
 "database/sql"

 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
)

const executionFields="id,request_id,operator_id,status,response_status,created_at,finished_at"
func scanExecution(row scanner)(entity.Execution,error){
 var e entity.Execution
 err:=row.Scan(&e.ID,&e.RequestID,&e.OperatorID,&e.Status,&e.ResponseStatus,&e.CreatedAt,&e.FinishedAt)
 return e,err
}
func(r *Repository) Environment(ctx context.Context,id string)(entity.Environment,error){
 var e entity.Environment
 err:=r.db.QueryRowContext(ctx,"SELECT id,code,name FROM internal_tool.environments WHERE id=$1",id).Scan(&e.ID,&e.Code,&e.Name)
 return e,err
}
func(r *Repository) ClaimExecution(ctx context.Context,actor entity.Operator,requestID,id string,version int64)(entity.Execution,error){
 tx,err:=r.db.BeginTx(ctx,nil);if err!=nil{return entity.Execution{},err};defer tx.Rollback()
 a,err:=scanApproval(tx.QueryRowContext(ctx,"SELECT "+approvalFields+" FROM internal_tool.approval_requests WHERE id=$1 FOR UPDATE",requestID));if err!=nil{return entity.Execution{},err}
 if actor.ID!=a.RequesterID&&actor.Persona!=entity.L2{return entity.Execution{},ErrForbidden}
 if a.Status!="APPROVED"||a.Version!=version{return entity.Execution{},ErrConflict}
 var active bool
 if err=tx.QueryRowContext(ctx,"SELECT EXISTS(SELECT 1 FROM internal_tool.operators WHERE id=$1 AND persona=$2 AND status='ACTIVE')",actor.ID,string(actor.Persona)).Scan(&active);err!=nil{return entity.Execution{},err}
 if !active{return entity.Execution{},ErrForbidden}
 e,err:=scanExecution(tx.QueryRowContext(ctx,"INSERT INTO internal_tool.executions(id,request_id,operator_id,status) VALUES($1,$2,$3,'RUNNING') ON CONFLICT(request_id) DO NOTHING RETURNING "+executionFields,id,requestID,actor.ID))
 if err==sql.ErrNoRows{return entity.Execution{},ErrConflict};if err!=nil{return e,err}
 if err=event(ctx,tx,actor.ID,requestID,"EXECUTION_STARTED");err!=nil{return e,err}
 return e,tx.Commit()
}
func(r *Repository) FinishExecution(ctx context.Context,e entity.Execution,status string,response *int)(entity.Execution,error){
 tx,err:=r.db.BeginTx(ctx,nil);if err!=nil{return e,err};defer tx.Rollback()
 e,err=scanExecution(tx.QueryRowContext(ctx,"UPDATE internal_tool.executions SET status=$2,response_status=$3,finished_at=now() WHERE id=$1 AND status='RUNNING' RETURNING "+executionFields,e.ID,status,response));if err!=nil{return e,err}
 if err=event(ctx,tx,e.OperatorID,e.RequestID,"EXECUTION_"+status);err!=nil{return e,err}
 return e,tx.Commit()
}
func(r *Repository) Execution(ctx context.Context,actor entity.Operator,requestID string)(entity.Execution,error){
 if _,err:=r.Get(ctx,actor,requestID);err!=nil{return entity.Execution{},err}
 return scanExecution(r.db.QueryRowContext(ctx,"SELECT "+executionFields+" FROM internal_tool.executions WHERE request_id=$1",requestID))
}
