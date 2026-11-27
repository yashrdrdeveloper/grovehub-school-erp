package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"strings"
	"grovehub/launcher/bootstrap"
	"grovehub/services/audit"
	"grovehub/services/infrastructure"
)

func main() {
	if err:=run();err!=nil {
		fmt.Fprintln(os.Stderr,"audit-replay:",err)
		os.Exit(1)
	}
}

func run() error {
	tenant:=flag.String("tenant-id","","Exact tenant identifier")
	event:=flag.String("event-id","","Exact original published event UUID")
	database:=flag.String("database","","Expected business database name")
	auditDatabase:=flag.String("audit-database","","Expected Audit database name; required with AUDIT_DATABASE_URL")
	operator:=flag.String("operator","","Operator label (required for apply)")
	apply:=flag.Bool("apply",false,"Persist the original audit event and replay receipt")
	flag.Parse()
	if flag.NArg()!=0 || *tenant=="" || *event=="" || *database=="" || (*apply && *operator=="") { return errors.New("tenant-id, event-id, database and operator for apply are required") }
	if strings.TrimSpace(os.Getenv("AUDIT_DATABASE_URL"))!="" && *auditDatabase=="" { return errors.New("audit-database is required with AUDIT_DATABASE_URL") }
	if *auditDatabase=="" { *auditDatabase=*database }
	ctx,cancel:=context.WithTimeout(context.Background(),15*time.Second)
	defer cancel()
	db,err:=bootstrap.NewPostgresDB(ctx)
	if err!=nil { return errors.New("could not connect to business database") }
	sqlDB,err:=db.DB()
	if err!=nil { return errors.New("could not access business database") }
	defer sqlDB.Close()
	var current string
	if err:=db.WithContext(ctx).Raw("SELECT current_database()").Scan(&current).Error;err!=nil || current!=*database { return errors.New("configured business database does not match --database") }
	auditDB,closeAudit,err:=bootstrap.NewAuditPostgresDB(ctx,db)
	if err!=nil { return err }
	defer closeAudit()
	if err:=auditDB.WithContext(ctx).Raw("SELECT current_database()").Scan(&current).Error;err!=nil || current!=*auditDatabase { return errors.New("configured Audit database does not match --audit-database") }
	result,err:=audit.Replay(ctx,auditDB,infrastructure.NewReplaySource(db),audit.ReplayRequest{TenantID:*tenant,EventID:*event,OperatorLabel:*operator,Apply:*apply})
	if err!=nil { return errors.New("replay failed; verify retained published event, tenant, schema and identity conflicts") }
	return json.NewEncoder(os.Stdout).Encode(result)
}
