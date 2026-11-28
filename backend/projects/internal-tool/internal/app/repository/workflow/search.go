package workflow

import (
 "context"
 "database/sql"

 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
)

func readPage[T any](rows *sql.Rows,p entity.PageRequest,scan func(scanner)(T,error))(entity.Page[T],error){
 defer rows.Close()
 records:=[]T{}
 for rows.Next(){v,err:=scan(rows);if err!=nil{return entity.Page[T]{},err};records=append(records,v)}
 if err:=rows.Err();err!=nil{return entity.Page[T]{},err}
 return entity.NewPage(p,records),nil
}
func(r *Repository) SearchEnvironments(ctx context.Context,p entity.PageRequest)(entity.Page[entity.Environment],error){
 rows,err:=r.db.QueryContext(ctx,"SELECT id,code,name FROM internal_tool.environments ORDER BY code LIMIT $1 OFFSET $2",p.RowSize+1,p.Offset())
 if err!=nil{return entity.Page[entity.Environment]{},err}
 return readPage(rows,p,func(row scanner)(entity.Environment,error){var e entity.Environment;err:=row.Scan(&e.ID,&e.Code,&e.Name);return e,err})
}
func(r *Repository) SearchTemplates(ctx context.Context,p entity.PageRequest)(entity.Page[entity.Template],error){
 rows,err:=r.db.QueryContext(ctx,"SELECT "+templateFields+" FROM internal_tool.api_templates ORDER BY code,version DESC LIMIT $1 OFFSET $2",p.RowSize+1,p.Offset())
 if err!=nil{return entity.Page[entity.Template]{},err};return readPage(rows,p,scanTemplate)
}
func(r *Repository) SearchApprovals(ctx context.Context,actor entity.Operator,p entity.PageRequest,inbox bool)(entity.Page[entity.Approval],error){
 var rows *sql.Rows;var err error
 if inbox{
  if actor.Persona!=entity.L2{return entity.Page[entity.Approval]{},ErrForbidden}
  rows,err=r.db.QueryContext(ctx,"SELECT "+approvalFields+" FROM internal_tool.approval_requests WHERE status='PENDING' ORDER BY created_at DESC,id DESC LIMIT $1 OFFSET $2",p.RowSize+1,p.Offset())
 }else{
  rows,err=r.db.QueryContext(ctx,"SELECT "+approvalFields+" FROM internal_tool.approval_requests WHERE requester_id=$1 ORDER BY created_at DESC,id DESC LIMIT $2 OFFSET $3",actor.ID,p.RowSize+1,p.Offset())
 }
 if err!=nil{return entity.Page[entity.Approval]{},err}
 return readPage(rows,p,scanApproval)
}
func(r *Repository) SearchEvents(ctx context.Context,actor entity.Operator,p entity.PageRequest)(entity.Page[entity.WorkflowEvent],error){
 if actor.Persona!=entity.L2{return entity.Page[entity.WorkflowEvent]{},ErrForbidden}
 rows,err:=r.db.QueryContext(ctx,"SELECT id,actor_id,resource_id,action,occurred_at,correlation_id FROM internal_tool.workflow_events ORDER BY id DESC LIMIT $1 OFFSET $2",p.RowSize+1,p.Offset())
 if err!=nil{return entity.Page[entity.WorkflowEvent]{},err}
 return readPage(rows,p,func(row scanner)(entity.WorkflowEvent,error){var e entity.WorkflowEvent;err:=row.Scan(&e.ID,&e.ActorID,&e.ResourceID,&e.Action,&e.OccurredAt,&e.CorrelationID);return e,err})
}
