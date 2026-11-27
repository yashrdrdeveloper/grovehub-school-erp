package operator

import (
 "context"
 "database/sql"
 "errors"

 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
)

var ErrForbidden=errors.New("operator management forbidden")
var ErrConflict=errors.New("operator management conflict")

type Summary struct{
 entity.Operator
 Status string `json:"status"`
}

func(r *Repository) List(ctx context.Context,actor entity.Operator)([]Summary,error){
 if actor.Persona!=entity.L2{return nil,ErrForbidden}
 rows,err:=r.db.QueryContext(ctx,"SELECT id,email,persona,status FROM internal_tool.operators ORDER BY email LIMIT 100");if err!=nil{return nil,err};defer rows.Close()
 result:=[]Summary{}
 for rows.Next(){var row Summary;if err=rows.Scan(&row.ID,&row.Email,&row.Persona,&row.Status);err!=nil{return nil,err};result=append(result,row)}
 return result,rows.Err()
}

// Serializing management changes prevents concurrent removal of the last L2.
func management(ctx context.Context,db *sql.DB,actor entity.Operator)(*sql.Tx,error){
 if actor.Persona!=entity.L2{return nil,ErrForbidden}
 tx,err:=db.BeginTx(ctx,nil);if err!=nil{return nil,err}
 if _,err=tx.ExecContext(ctx,"SELECT pg_advisory_xact_lock(80228102)");err!=nil{tx.Rollback();return nil,err}
 var allowed bool
 if err=tx.QueryRowContext(ctx,"SELECT EXISTS(SELECT 1 FROM internal_tool.operators WHERE id=$1 AND persona='L2' AND status='ACTIVE')",actor.ID).Scan(&allowed);err!=nil||!allowed{tx.Rollback();if err!=nil{return nil,err};return nil,ErrForbidden}
 return tx,nil
}
func(r *Repository) Create(ctx context.Context,actor entity.Operator,c entity.Credential)error{
 tx,err:=management(ctx,r.db,actor);if err!=nil{return err};defer tx.Rollback()
 result,err:=tx.ExecContext(ctx,"INSERT INTO internal_tool.operators(id,email,persona,password_salt,password_hash) VALUES($1,$2,$3,$4,$5) ON CONFLICT(email) DO NOTHING",c.Operator.ID,c.Operator.Email,c.Operator.Persona,c.Salt,c.Hash)
 if err!=nil{return err};n,err:=result.RowsAffected();if err!=nil{return err};if n!=1{return ErrConflict}
 if _,err=tx.ExecContext(ctx,"INSERT INTO internal_tool.workflow_events(actor_id,resource_id,action) VALUES($1,$2,'OPERATOR_CREATED')",actor.ID,c.Operator.ID);err!=nil{return err}
 return tx.Commit()
}
func(r *Repository) SetStatus(ctx context.Context,actor entity.Operator,id,status string)error{
 if actor.ID==id{return ErrForbidden}
 tx,err:=management(ctx,r.db,actor);if err!=nil{return err};defer tx.Rollback()
 var persona,current string
 if err=tx.QueryRowContext(ctx,"SELECT persona,status FROM internal_tool.operators WHERE id=$1 FOR UPDATE",id).Scan(&persona,&current);err!=nil{return err}
 if current==status{return ErrConflict}
 if persona=="L2"&&status=="DISABLED"{
  var count int
  if err=tx.QueryRowContext(ctx,"SELECT count(*) FROM internal_tool.operators WHERE persona='L2' AND status='ACTIVE'").Scan(&count);err!=nil{return err}
  if count<=1{return ErrConflict}
 }
 if _,err=tx.ExecContext(ctx,"UPDATE internal_tool.operators SET status=$2 WHERE id=$1",id,status);err!=nil{return err}
 if _,err=tx.ExecContext(ctx,"UPDATE internal_tool.sessions SET revoked_at=now() WHERE operator_id=$1 AND revoked_at IS NULL",id);err!=nil{return err}
 if _,err=tx.ExecContext(ctx,"INSERT INTO internal_tool.workflow_events(actor_id,resource_id,action) VALUES($1,$2,$3)",actor.ID,id,"OPERATOR_"+status);err!=nil{return err}
 return tx.Commit()
}
