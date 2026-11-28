package workflow

import (
 "context"
 corehttp "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/repository/activity"
 "database/sql"
 "encoding/json"
 "errors"

 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
)

var ErrConflict=errors.New("workflow conflict")
var ErrForbidden=errors.New("workflow forbidden")
type Repository struct{db *sql.DB}
func New(db *sql.DB)*Repository{return &Repository{db}}

func event(ctx context.Context,tx *sql.Tx,actor,id,action string)error{
 return activity.Record(ctx,tx,actor,id,action)
}
func(r *Repository) CreateEnvironment(ctx context.Context,actor string,e entity.Environment)error{
 tx,err:=r.db.BeginTx(ctx,nil);if err!=nil{return err};defer tx.Rollback()
 result,err:=tx.ExecContext(ctx,"INSERT INTO internal_tool.environments(id,code,name,created_by) VALUES($1,$2,$3,$4) ON CONFLICT(code) DO NOTHING",e.ID,e.Code,e.Name,actor);if err!=nil{return err}
 n,err:=result.RowsAffected();if err!=nil{return err};if n!=1{return ErrConflict}
 if err=event(ctx,tx,actor,e.ID,"ENVIRONMENT_CREATED");err!=nil{return err};return tx.Commit()
}
func(r *Repository) Environments(ctx context.Context)([]entity.Environment,error){
 rows,err:=r.db.QueryContext(ctx,"SELECT id,code,name FROM internal_tool.environments ORDER BY code LIMIT 100");if err!=nil{return nil,err};defer rows.Close()
 result:=[]entity.Environment{};for rows.Next(){var e entity.Environment;if err:=rows.Scan(&e.ID,&e.Code,&e.Name);err!=nil{return nil,err};result=append(result,e)};return result,rows.Err()
}
func(r *Repository) CreateTemplate(ctx context.Context,actor string,t *entity.Template)error{
 tx,err:=r.db.BeginTx(ctx,nil);if err!=nil{return err};defer tx.Rollback()
 if _,err=tx.ExecContext(ctx,"SELECT pg_advisory_xact_lock(80228103)");err!=nil{return err}
 if err=tx.QueryRowContext(ctx,"SELECT COALESCE(MAX(version),0)+1 FROM internal_tool.api_templates WHERE code=$1",t.Code).Scan(&t.Version);err!=nil{return err}
 headers,err:=json.Marshal(t.Headers);if err!=nil{return err};schema,err:=json.Marshal(t.BodySchema);if err!=nil{return err}
 if _,err=tx.ExecContext(ctx,"INSERT INTO internal_tool.api_templates(id,code,version,module,method,endpoint,scope,headers,body_schema,created_by) VALUES($1,$2,$3,$4,$5,$6,$7,$8::jsonb,$9::jsonb,$10)",t.ID,t.Code,t.Version,t.Module,t.Method,t.Endpoint,t.Scope,string(headers),string(schema),actor);err!=nil{return err}
 if err=event(ctx,tx,actor,t.ID,"TEMPLATE_CREATED");err!=nil{return err};return tx.Commit()
}
type scanner interface{Scan(...any)error}
const templateFields="id,code,version,module,method,endpoint,scope,headers,body_schema"
func scanTemplate(row scanner)(entity.Template,error){
 var t entity.Template;var headers,schema []byte
 if err:=row.Scan(&t.ID,&t.Code,&t.Version,&t.Module,&t.Method,&t.Endpoint,&t.Scope,&headers,&schema);err!=nil{return t,err}
 if err:=json.Unmarshal(headers,&t.Headers);err!=nil{return t,err};return t,json.Unmarshal(schema,&t.BodySchema)
}
func(r *Repository) Template(ctx context.Context,id string)(entity.Template,error){
 return scanTemplate(r.db.QueryRowContext(ctx,"SELECT "+templateFields+" FROM internal_tool.api_templates WHERE id=$1",id))
}
func(r *Repository) Templates(ctx context.Context)([]entity.Template,error){
 rows,err:=r.db.QueryContext(ctx,"SELECT "+templateFields+" FROM internal_tool.api_templates ORDER BY code,version DESC LIMIT 100");if err!=nil{return nil,err};defer rows.Close()
 result:=[]entity.Template{};for rows.Next(){t,err:=scanTemplate(rows);if err!=nil{return nil,err};result=append(result,t)};return result,rows.Err()
}
func(r *Repository) Submit(ctx context.Context,a entity.Approval,auto bool) (entity.Approval,error){
 ctx=corehttp.Ensure(ctx)
 tx,err:=r.db.BeginTx(ctx,nil);if err!=nil{return a,err};defer tx.Rollback()
 var exists bool
 if err=tx.QueryRowContext(ctx,"SELECT EXISTS(SELECT 1 FROM internal_tool.environments WHERE id=$1)",a.EnvironmentID).Scan(&exists);err!=nil{return a,err};if !exists{return a,sql.ErrNoRows}
 status:="PENDING";if auto{status="APPROVED"}
 query:="INSERT INTO internal_tool.approval_requests(id,template_id,environment_id,school_id,input,reason,requester_id,status,reviewer_id,decision_comment,decided_at) VALUES($1,$2,$3,$4,$5::jsonb,$6,$7,$8,CASE WHEN $8='APPROVED' THEN $7::uuid END,CASE WHEN $8='APPROVED' THEN 'L2 self approval' END,CASE WHEN $8='APPROVED' THEN now() END) RETURNING "+approvalFields
 a,err=scanApproval(tx.QueryRowContext(ctx,query,a.ID,a.TemplateID,a.EnvironmentID,a.SchoolID,string(a.Input),a.Reason,a.RequesterID,status));if err!=nil{return a,err}
 if err=event(ctx,tx,a.RequesterID,a.ID,"REQUEST_SUBMITTED");err!=nil{return a,err}
 if auto{if err=event(ctx,tx,a.RequesterID,a.ID,"REQUEST_AUTO_APPROVED");err!=nil{return a,err}}
 return a,tx.Commit()
}
const approvalFields="id,template_id,environment_id,school_id,input,reason,requester_id,status,version,reviewer_id,decision_comment,created_at,decided_at"
func scanApproval(row scanner)(entity.Approval,error){
 var a entity.Approval
 err:=row.Scan(&a.ID,&a.TemplateID,&a.EnvironmentID,&a.SchoolID,&a.Input,&a.Reason,&a.RequesterID,&a.Status,&a.Version,&a.ReviewerID,&a.DecisionComment,&a.CreatedAt,&a.DecidedAt)
 return a,err
}
func(r *Repository) List(ctx context.Context,actor entity.Operator,inbox bool)([]entity.Approval,error){
 query:="SELECT "+approvalFields+" FROM internal_tool.approval_requests WHERE requester_id=$1 ORDER BY created_at DESC,id DESC LIMIT 100"
 arg:=[]any{actor.ID}
 if inbox {if actor.Persona!=entity.L2{return nil,ErrForbidden};query="SELECT "+approvalFields+" FROM internal_tool.approval_requests WHERE status='PENDING' ORDER BY created_at DESC,id DESC LIMIT 100";arg=nil}
 rows,err:=r.db.QueryContext(ctx,query,arg...);if err!=nil{return nil,err};defer rows.Close()
 result:=[]entity.Approval{};for rows.Next(){a,err:=scanApproval(rows);if err!=nil{return nil,err};result=append(result,a)};return result,rows.Err()
}
func(r *Repository) Get(ctx context.Context,actor entity.Operator,id string)(entity.Approval,error){
 return scanApproval(r.db.QueryRowContext(ctx,"SELECT "+approvalFields+" FROM internal_tool.approval_requests WHERE id=$1 AND (requester_id=$2 OR $3='L2')",id,actor.ID,string(actor.Persona)))
}
func(r *Repository) Decide(ctx context.Context,actor entity.Operator,id string,version int64,status,comment string)(entity.Approval,error){
 tx,err:=r.db.BeginTx(ctx,nil);if err!=nil{return entity.Approval{},err};defer tx.Rollback()
 a,err:=decide(ctx,tx,actor,id,version,status,comment)
 if err!=nil{return a,err}
 return a,tx.Commit()
}

func decide(ctx context.Context,tx *sql.Tx,actor entity.Operator,id string,version int64,status,comment string)(entity.Approval,error){
 a,err:=scanApproval(tx.QueryRowContext(ctx,"SELECT "+approvalFields+" FROM internal_tool.approval_requests WHERE id=$1 FOR UPDATE",id));if err!=nil{return a,err}
 if status=="CANCELLED" {if actor.ID!=a.RequesterID{return entity.Approval{},ErrForbidden}} else if actor.Persona!=entity.L2{return entity.Approval{},ErrForbidden}
 if a.Status!="PENDING"||a.Version!=version{return entity.Approval{},ErrConflict}
 var active bool
 if err=tx.QueryRowContext(ctx,"SELECT EXISTS(SELECT 1 FROM internal_tool.operators WHERE id=$1 AND persona=$2 AND status='ACTIVE')",actor.ID,string(actor.Persona)).Scan(&active);err!=nil{return entity.Approval{},err}
 if !active{return entity.Approval{},ErrForbidden}
 if status=="APPROVED"&&actor.ID==a.RequesterID&&a.SchoolID!=nil{
  var method string
  if err=tx.QueryRowContext(ctx,"SELECT method FROM internal_tool.api_templates WHERE id=$1",a.TemplateID).Scan(&method);err!=nil{return entity.Approval{},err}
  if method!="GET"{return entity.Approval{},ErrForbidden}
 }
 a,err=scanApproval(tx.QueryRowContext(ctx,"UPDATE internal_tool.approval_requests SET status=$2,version=version+1,reviewer_id=$3,decision_comment=$4,decided_at=now() WHERE id=$1 RETURNING "+approvalFields,id,status,actor.ID,comment));if err!=nil{return a,err}
 if err=event(ctx,tx,actor.ID,id,"REQUEST_"+status);err!=nil{return a,err};return a,nil
}
