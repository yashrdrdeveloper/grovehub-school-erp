package repository

import (
 "context"
 "database/sql"
 "errors"

 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/registry/internal/app/entity"
)

var ErrConflict=errors.New("school registry conflict")

func Register(ctx context.Context,db *sql.DB,s entity.School)error{
 tx,err:=db.BeginTx(ctx,nil);if err!=nil{return err};defer tx.Rollback()
 if _,err=tx.ExecContext(ctx,"INSERT INTO school_registry.schools(id,code,name) VALUES($1,$2,$3) ON CONFLICT DO NOTHING",s.ID,s.Code,s.Name);err!=nil{return err}
 var same bool
 if err=tx.QueryRowContext(ctx,"SELECT EXISTS(SELECT 1 FROM school_registry.schools WHERE id=$1 AND code=$2 AND name=$3)",s.ID,s.Code,s.Name).Scan(&same);err!=nil{return err}
 if !same{return ErrConflict}
 return tx.Commit()
}

func Lookup(ctx context.Context,db *sql.DB,ids []string)([]entity.School,error){
 // Bound by the caller to 100 canonical UUIDs; no unrestricted listing.
 rows,err:=db.QueryContext(ctx,"SELECT id,code,name,status FROM school_registry.schools WHERE id=ANY($1::uuid[]) ORDER BY id",ids)
 if err!=nil{return nil,err};defer rows.Close()
 result:=[]entity.School{}
 for rows.Next(){var school entity.School;if err=rows.Scan(&school.ID,&school.Code,&school.Name,&school.Status);err!=nil{return nil,err};result=append(result,school)}
 return result,rows.Err()
}
