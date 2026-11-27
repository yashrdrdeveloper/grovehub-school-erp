package session

import (
 "context"
 "database/sql"
 "time"

 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
)

type Repository struct{db *sql.DB}
func New(db *sql.DB)*Repository{return &Repository{db:db}}

func(r *Repository) Create(ctx context.Context,operatorID string,token,csrf []byte,expires time.Time)error{
 _,err:=r.db.ExecContext(ctx,"INSERT INTO internal_tool.sessions(token_hash,operator_id,csrf_hash,expires_at) VALUES($1,$2,$3,$4)",token,operatorID,csrf,expires)
 return err
}

func(r *Repository) Verify(ctx context.Context,token []byte)(entity.Operator,[]byte,error){
 var operator entity.Operator;var csrf []byte
 err:=r.db.QueryRowContext(ctx,`SELECT o.id,o.email,o.persona,s.csrf_hash FROM internal_tool.sessions s JOIN internal_tool.operators o ON o.id=s.operator_id WHERE s.token_hash=$1 AND s.revoked_at IS NULL AND s.expires_at>now() AND o.status='ACTIVE'`,token).Scan(&operator.ID,&operator.Email,&operator.Persona,&csrf)
 return operator,csrf,err
}

func(r *Repository) Revoke(ctx context.Context,token []byte)error{
 _,err:=r.db.ExecContext(ctx,"UPDATE internal_tool.sessions SET revoked_at=now() WHERE token_hash=$1 AND revoked_at IS NULL",token);return err
}

// ClaimLoginAttempt uses a persistent, atomic per-identifier budget. Attempts
// are consumed before password hashing, including for unknown accounts.
func(r *Repository) ClaimLoginAttempt(ctx context.Context,key []byte)(bool,error){
 var attempts int
 err:=r.db.QueryRowContext(ctx,`INSERT INTO internal_tool.login_attempts(email_hash,window_start,attempts) VALUES($1,now(),1)
 ON CONFLICT(email_hash) DO UPDATE SET
 attempts=CASE WHEN internal_tool.login_attempts.window_start<=now()-interval '5 minutes' THEN 1 ELSE LEAST(internal_tool.login_attempts.attempts+1,6) END,
 window_start=CASE WHEN internal_tool.login_attempts.window_start<=now()-interval '5 minutes' THEN now() ELSE internal_tool.login_attempts.window_start END
 RETURNING attempts`,key).Scan(&attempts)
 return attempts<=5,err
}
