package operator

import (
 "context"
 "database/sql"
 "errors"

 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
)

var ErrBootstrapClosed=errors.New("operator bootstrap is already complete")
type Repository struct{db *sql.DB}
func New(db *sql.DB)*Repository{return &Repository{db:db}}

func(r *Repository) Credential(ctx context.Context,email string)(entity.Credential,error){
 var c entity.Credential
 err:=r.db.QueryRowContext(ctx,"SELECT id,email,persona,password_salt,password_hash,status='ACTIVE' FROM internal_tool.operators WHERE email=$1",email).Scan(&c.Operator.ID,&c.Operator.Email,&c.Operator.Persona,&c.Salt,&c.Hash,&c.Active)
 return c,err
}

// Bootstrap is serialized and only succeeds when no operator exists.
func(r *Repository) Bootstrap(ctx context.Context,c entity.Credential)error{
 tx,err:=r.db.BeginTx(ctx,nil);if err!=nil{return err};defer tx.Rollback()
 if _,err=tx.ExecContext(ctx,"SELECT pg_advisory_xact_lock(80228102)");err!=nil{return err}
 var exists bool
 if err=tx.QueryRowContext(ctx,"SELECT EXISTS(SELECT 1 FROM internal_tool.operators)").Scan(&exists);err!=nil{return err}
 if exists{return ErrBootstrapClosed}
 if _,err=tx.ExecContext(ctx,"INSERT INTO internal_tool.operators(id,email,persona,password_salt,password_hash) VALUES($1,$2,$3,$4,$5)",c.Operator.ID,c.Operator.Email,c.Operator.Persona,c.Salt,c.Hash);err!=nil{return err}
 return tx.Commit()
}
