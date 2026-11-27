package repository

import (
 "context"
 "database/sql"
 "errors"
)

var ErrInactive=errors.New("school inactive")

// Rename commits the name and immutable receipt together. Identical retries
// acknowledge the original execution without changing the current school name.
func Rename(ctx context.Context,db *sql.DB,school,execution,subject,expected,name string,version int64,correlationID *string)error{
 tx,err:=db.BeginTx(ctx,nil);if err!=nil{return err};defer tx.Rollback()
 var current,status string
 var currentVersion int64
 if err=tx.QueryRowContext(ctx,"SELECT name,status,version FROM school_registry.schools WHERE id=$1 FOR UPDATE",school).Scan(&current,&status,&currentVersion);err!=nil{return err}
 var oldSchool,oldSubject,oldExpected,oldName string
 var oldVersion sql.NullInt64
 var sameCorrelation bool
 err=tx.QueryRowContext(ctx,"SELECT school_id,workload_subject,previous_name,new_name,previous_version, correlation_id IS NOT DISTINCT FROM CAST($2 AS uuid) FROM school_registry.rename_events WHERE execution_id=$1",execution,correlationID).Scan(&oldSchool,&oldSubject,&oldExpected,&oldName,&oldVersion,&sameCorrelation)
 if err==nil{
  if oldSchool!=school||oldSubject!=subject||oldExpected!=expected||oldName!=name||!oldVersion.Valid||oldVersion.Int64!=version||!sameCorrelation{return ErrConflict}
  return tx.Commit()
 }
 if !errors.Is(err,sql.ErrNoRows){return err}
 if status!="ACTIVE"{return ErrInactive}
 if currentVersion!=version||current!=expected||current==name{return ErrConflict}
 result,err:=tx.ExecContext(ctx,"INSERT INTO school_registry.rename_events(school_id,execution_id,workload_subject,previous_name,new_name,previous_version,correlation_id) VALUES($1,$2,$3,$4,$5,$6,$7) ON CONFLICT(execution_id) DO NOTHING",school,execution,subject,expected,name,version,correlationID)
 if err!=nil{return err}
 count,err:=result.RowsAffected();if err!=nil{return err};if count!=1{return ErrConflict}
 if _,err=tx.ExecContext(ctx,"UPDATE school_registry.schools SET name=$2 WHERE id=$1",school,name);err!=nil{return err}
 return tx.Commit()
}
