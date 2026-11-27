package migrations

import (
 "context"
 "database/sql"
 _ "embed"
 "errors"
)

//go:embed 001_schools.sql
var schema string

//go:embed 002_school_rename.sql
var renameSchema string

//go:embed 003_school_version.sql
var versionSchema string

//go:embed 004_rename_delivery.sql
var deliverySchema string

const Version=4

func Apply(ctx context.Context,db *sql.DB)error{
 tx,err:=db.BeginTx(ctx,nil);if err!=nil{return err};defer tx.Rollback()
 if _,err=tx.ExecContext(ctx,"SELECT pg_advisory_xact_lock(80228110)");err!=nil{return err}
 if _,err=tx.ExecContext(ctx,"CREATE SCHEMA IF NOT EXISTS school_registry; CREATE TABLE IF NOT EXISTS school_registry.schema_version(version integer PRIMARY KEY)");err!=nil{return err}
 var version int
 if err=tx.QueryRowContext(ctx,"SELECT COALESCE(MAX(version),0) FROM school_registry.schema_version").Scan(&version);err!=nil{return err}
 if version>Version{return errors.New("school registry schema is newer than binary")}
 if version==0{
  if _,err=tx.ExecContext(ctx,schema);err!=nil{return err}
  if _,err=tx.ExecContext(ctx,"INSERT INTO school_registry.schema_version(version) VALUES(1)");err!=nil{return err}
 }
 if version<2{
  if _,err=tx.ExecContext(ctx,renameSchema);err!=nil{return err}
  if _,err=tx.ExecContext(ctx,"INSERT INTO school_registry.schema_version(version) VALUES(2)");err!=nil{return err}
 }
 if version<3{
  if _,err=tx.ExecContext(ctx,versionSchema);err!=nil{return err}
  if _,err=tx.ExecContext(ctx,"INSERT INTO school_registry.schema_version(version) VALUES(3)");err!=nil{return err}
 }
 if version<4{
  if _,err=tx.ExecContext(ctx,deliverySchema);err!=nil{return err}
  if _,err=tx.ExecContext(ctx,"INSERT INTO school_registry.schema_version(version) VALUES(4)");err!=nil{return err}
 }
 return tx.Commit()
}
func Check(ctx context.Context,db *sql.DB)error{
 var version int
 if err:=db.QueryRowContext(ctx,"SELECT COALESCE(MAX(version),0) FROM school_registry.schema_version").Scan(&version);err!=nil||version!=Version{return errors.New("school registry migration required")}
 return nil
}
