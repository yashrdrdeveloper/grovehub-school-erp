package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/postgres"
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
	database:=flag.String("database","","Expected database name")
	operator:=flag.String("operator","","Operator label (required for apply)")
	apply:=flag.Bool("apply",false,"Persist the original audit event and replay receipt")
	flag.Parse()
	if flag.NArg()!=0 || *tenant=="" || *event=="" || *database=="" || (*apply && *operator=="") { return errors.New("tenant-id, event-id, database and operator for apply are required") }
	dsn:=os.Getenv("DATABASE_URL")
	if dsn=="" { return errors.New("DATABASE_URL is required") }
	db,err:=postgres.NewDB(dsn)
	if err!=nil { return errors.New("could not connect to database") }
	sqlDB,err:=db.DB()
	if err!=nil { return errors.New("could not access database") }
	defer sqlDB.Close()
	ctx,cancel:=context.WithTimeout(context.Background(),15*time.Second)
	defer cancel()
	var current string
	if err:=db.WithContext(ctx).Raw("SELECT current_database()").Scan(&current).Error;err!=nil || current!=*database { return errors.New("configured database does not match --database") }
	result,err:=audit.Replay(ctx,db,infrastructure.NewReplaySource(db),audit.ReplayRequest{TenantID:*tenant,EventID:*event,OperatorLabel:*operator,Apply:*apply})
	if err!=nil { return errors.New("replay failed; verify retained published event, tenant, schema and identity conflicts") }
	return json.NewEncoder(os.Stdout).Encode(result)
}
