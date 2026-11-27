package repository

import (
 "context"
 "database/sql"
 "errors"
)

var ErrInactive=errors.New("school inactive")

// Rename commits the name and immutable receipt together. Identical retries
// acknowledge the original execution without changing the current school name.
func Rename(ctx context.Context,db *sql.DB,school,execution,subject,expected,name string)error{
 tx,err:=db.BeginTx(ctx,nil);if err!=nil{return err};defer tx.Rollback()
 var current,status string
 if err=tx.QueryRowContext(ctx,"SELECT name,status FROM school_registry.schools WHERE id=$1 FOR UPDATE",school).Scan(&current,&status);err!=nil{return err}
 var oldSchool,oldSubject,oldExpected,oldName string
 err=tx.QueryRowContext(ctx,"SELECT school_id,workload_subject,previous_name,new_name FROM school_registry.rename_events WHERE execution_id=$1",execution).Scan(&oldSchool,&oldSubject,&oldExpected,&oldName)
 if err==nil{
  if oldSchool!=school||oldSubject!=subject||oldExpected!=expected||oldName!=name{return ErrConflict}
  return tx.Commit()
 }
 if !errors.Is(err,sql.ErrNoRows){return err}
 if status!="ACTIVE"{return ErrInactive}
 if current!=expected||current==name{return ErrConflict}
 result,err:=tx.ExecContext(ctx,"INSERT INTO school_registry.rename_events(school_id,execution_id,workload_subject,previous_name,new_name) VALUES($1,$2,$3,$4,$5) ON CONFLICT(execution_id) DO NOTHING",school,execution,subject,expected,name)
 if err!=nil{return err}
 count,err:=result.RowsAffected();if err!=nil{return err};if count!=1{return ErrConflict}
 if _,err=tx.ExecContext(ctx,"UPDATE school_registry.schools SET name=$2 WHERE id=$1",school,name);err!=nil{return err}
 return tx.Commit()
}
