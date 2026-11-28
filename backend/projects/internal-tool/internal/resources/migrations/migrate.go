package migrations

import (
 "context"
 "database/sql"
 _ "embed"
 "errors"
)

//go:embed 001_operator_identity.sql
var operatorIdentity string

//go:embed 002_approval_registry.sql
var approvalRegistry string

//go:embed 003_execution.sql
var executionSchema string

//go:embed 004_execution_reconciliation.sql
var reconciliationSchema string

//go:embed 005_audit_outbox.sql
var auditOutboxSchema string

//go:embed 006_authentication_audit.sql
var authenticationAuditSchema string

const Version = 6

// Apply is an explicit operator command, never an implicit startup migration.
func Apply(ctx context.Context, db *sql.DB) error {
 tx,err:=db.BeginTx(ctx,nil)
 if err!=nil{return err}
 defer tx.Rollback()
 if _,err=tx.ExecContext(ctx,"SELECT pg_advisory_xact_lock(80228101)");err!=nil{return err}
 if _,err=tx.ExecContext(ctx,"CREATE SCHEMA IF NOT EXISTS internal_tool; CREATE TABLE IF NOT EXISTS internal_tool.schema_version (version integer PRIMARY KEY, applied_at timestamptz NOT NULL DEFAULT now())");err!=nil{return err}
 var version int
 if err=tx.QueryRowContext(ctx,"SELECT COALESCE(MAX(version),0) FROM internal_tool.schema_version").Scan(&version);err!=nil{return err}
 if version>Version{return errors.New("database schema is newer than this binary")}
 if version==0 {
  if _,err=tx.ExecContext(ctx,operatorIdentity);err!=nil{return err}
  if _,err=tx.ExecContext(ctx,"INSERT INTO internal_tool.schema_version(version) VALUES (1)");err!=nil{return err}
 }
 if version<2 {
  if _,err=tx.ExecContext(ctx,approvalRegistry);err!=nil{return err}
  if _,err=tx.ExecContext(ctx,"INSERT INTO internal_tool.schema_version(version) VALUES (2)");err!=nil{return err}
 }
 if version<3 {
  if _,err=tx.ExecContext(ctx,executionSchema);err!=nil{return err}
  if _,err=tx.ExecContext(ctx,"INSERT INTO internal_tool.schema_version(version) VALUES (3)");err!=nil{return err}
 }
 if version<4 {
  if _,err=tx.ExecContext(ctx,reconciliationSchema);err!=nil{return err}
  if _,err=tx.ExecContext(ctx,"INSERT INTO internal_tool.schema_version(version) VALUES (4)");err!=nil{return err}
 }
 if version<5 {
  if _,err=tx.ExecContext(ctx,auditOutboxSchema);err!=nil{return err}
  if _,err=tx.ExecContext(ctx,"INSERT INTO internal_tool.schema_version(version) VALUES (5)");err!=nil{return err}
 }
 if version<6 {
  if _,err=tx.ExecContext(ctx,authenticationAuditSchema);err!=nil{return err}
  if _,err=tx.ExecContext(ctx,"INSERT INTO internal_tool.schema_version(version) VALUES (6)");err!=nil{return err}
 }
 return tx.Commit()
}

func Check(ctx context.Context,db *sql.DB)error{
 var version int
 if err:=db.QueryRowContext(ctx,"SELECT COALESCE(MAX(version),0) FROM internal_tool.schema_version").Scan(&version);err!=nil||version!=Version{return errors.New("internal tool database migration required")}
 return nil
}
