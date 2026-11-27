package session

import (
 "context"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/repository/activity"
 "database/sql"
 "time"

 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
)

type Repository struct{db *sql.DB}
func New(db *sql.DB)*Repository{return &Repository{db:db}}

func(r *Repository) Create(ctx context.Context,operatorID string,token,csrf []byte,expires time.Time)error{
 tx,err:=r.db.BeginTx(ctx,nil);if err!=nil{return err};defer tx.Rollback()
 // Serialize against operator disable, which revokes sessions in the same transaction.
 var status string
 if err=tx.QueryRowContext(ctx,"SELECT status FROM internal_tool.operators WHERE id=$1 FOR UPDATE",operatorID).Scan(&status);err!=nil{return err}
 if status!="ACTIVE"{return sql.ErrNoRows}
 if _,err=tx.ExecContext(ctx,"INSERT INTO internal_tool.sessions(token_hash,operator_id,csrf_hash,expires_at) VALUES($1,$2,$3,$4)",token,operatorID,csrf,expires);err!=nil{return err}
 if err=activity.Record(ctx,tx,operatorID,operatorID,"AUTH_LOGIN_SUCCEEDED");err!=nil{return err}
 return tx.Commit()
}

func(r *Repository) Verify(ctx context.Context,token []byte)(entity.Operator,[]byte,error){
 var operator entity.Operator;var csrf []byte
 err:=r.db.QueryRowContext(ctx,`SELECT o.id,o.email,o.persona,s.csrf_hash FROM internal_tool.sessions s JOIN internal_tool.operators o ON o.id=s.operator_id WHERE s.token_hash=$1 AND s.revoked_at IS NULL AND s.expires_at>now() AND o.status='ACTIVE'`,token).Scan(&operator.ID,&operator.Email,&operator.Persona,&csrf)
 return operator,csrf,err
}

func(r *Repository) Revoke(ctx context.Context,token []byte)error{
 tx,err:=r.db.BeginTx(ctx,nil);if err!=nil{return err};defer tx.Rollback()
 var operator string
 err=tx.QueryRowContext(ctx,"UPDATE internal_tool.sessions SET revoked_at=now() WHERE token_hash=$1 AND revoked_at IS NULL RETURNING operator_id",token).Scan(&operator)
 if err==sql.ErrNoRows{return nil};if err!=nil{return err}
 if err=activity.Record(ctx,tx,operator,operator,"AUTH_LOGOUT_SUCCEEDED");err!=nil{return err}
 return tx.Commit()
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

func(r *Repository) CountAuthenticationOutcome(ctx context.Context,outcome string)error{
 _,err:=r.db.ExecContext(ctx,"INSERT INTO internal_tool.authentication_counters(bucket,outcome,attempts) VALUES(date_trunc('hour',now()),$1,1) ON CONFLICT(bucket,outcome) DO UPDATE SET attempts=internal_tool.authentication_counters.attempts+1",outcome)
 return err
}
