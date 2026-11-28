package repository

import (
 "context"
 "database/sql"
 "errors"
)

// A shared execution lock prevents reusing an execution UUID across mutation
// types or schools. Hash collisions only serialize unrelated transactions.
func lockExecution(ctx context.Context,tx *sql.Tx,id string)error{
 _,err:=tx.ExecContext(ctx,"SELECT pg_advisory_xact_lock(hashtextextended($1,80228111))",id)
 return err
}
func ChangeStatus(ctx context.Context,db *sql.DB,school,execution,subject,expected,status string,version int64,correlationID *string)error{
 tx,err:=db.BeginTx(ctx,nil);if err!=nil{return err};defer tx.Rollback()
 if err=lockExecution(ctx,tx,execution);err!=nil{return err}
 var occupied bool
 if err=tx.QueryRowContext(ctx,"SELECT EXISTS(SELECT 1 FROM school_registry.rename_events WHERE execution_id=$1)",execution).Scan(&occupied);err!=nil{return err}
 if occupied{return ErrConflict}
 var current string;var currentVersion int64
 if err=tx.QueryRowContext(ctx,"SELECT status,version FROM school_registry.schools WHERE id=$1 FOR UPDATE",school).Scan(&current,&currentVersion);err!=nil{return err}
 var oldSchool,oldSubject,oldExpected,oldStatus string;var oldVersion int64;var sameCorrelation bool
 err=tx.QueryRowContext(ctx,"SELECT school_id,workload_subject,previous_status,new_status,previous_version,correlation_id IS NOT DISTINCT FROM CAST($2 AS uuid) FROM school_registry.status_events WHERE execution_id=$1",execution,correlationID).Scan(&oldSchool,&oldSubject,&oldExpected,&oldStatus,&oldVersion,&sameCorrelation)
 if err==nil{
  if oldSchool!=school||oldSubject!=subject||oldExpected!=expected||oldStatus!=status||oldVersion!=version||!sameCorrelation{return ErrConflict}
  return tx.Commit()
 }
 if !errors.Is(err,sql.ErrNoRows){return err}
 if currentVersion!=version||current!=expected||current==status{return ErrConflict}
 if _,err=tx.ExecContext(ctx,"INSERT INTO school_registry.status_events(school_id,execution_id,workload_subject,previous_status,new_status,previous_version,correlation_id) VALUES($1,$2,$3,$4,$5,$6,$7)",school,execution,subject,expected,status,version,correlationID);err!=nil{return err}
 if _,err=tx.ExecContext(ctx,"UPDATE school_registry.schools SET status=$2 WHERE id=$1",school,status);err!=nil{return err}
 return tx.Commit()
}
