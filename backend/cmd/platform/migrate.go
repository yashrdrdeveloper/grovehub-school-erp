package main

import (
    "context"
    "fmt"
    "os"
    "time"
    "gorm.io/gorm"

    "grovehub/launcher/bootstrap"
    "grovehub/services/accounting"
    "grovehub/services/foundation"
    "grovehub/services/audit"
    "grovehub/services/identity"
    "grovehub/services/infrastructure"
    "grovehub/services/integration"
)

// runMigrationCommand applies embedded service migrations before app startup.
func runMigrationCommand() int {
    if len(os.Args)!=3 || os.Args[1]!="migrate" {
        _,_=fmt.Fprintln(os.Stderr,"usage: groovehub-platform migrate foundation|infrastructure|integration|identity|audit|accounting|all")
        return 2
    }
    service:=os.Args[2]
    if service!="foundation" && service!="infrastructure" && service!="integration" && service!="identity" && service!="audit" && service!="accounting" && service!="all" {
        _,_=fmt.Fprintln(os.Stderr,"usage: groovehub-platform migrate foundation|infrastructure|integration|identity|audit|accounting|all")
        return 2
    }
    ctx,cancel:=context.WithTimeout(context.Background(),10*time.Minute)
    defer cancel()
    if service=="audit" { return migrateAudit(ctx,nil) }
    db,err:=bootstrap.NewPostgresDB(ctx)
    if err!=nil{_,_=fmt.Fprintf(os.Stderr,"open migration database: %v\n",err);return 1}
    sqlDB,err:=db.DB()
    if err!=nil{_,_=fmt.Fprintf(os.Stderr,"access migration database: %v\n",err);return 1}
    defer func(){_=sqlDB.Close()}()
    if service=="foundation" || service=="all" {
        if err:=foundation.Migrate(ctx,db);err!=nil{
            _,_=fmt.Fprintf(os.Stderr,"migrate foundation: %v\n",err)
            return 1
        }
    }
    if service=="infrastructure" || service=="all" {
        if err:=infrastructure.Migrate(ctx,db);err!=nil{
            _,_=fmt.Fprintf(os.Stderr,"migrate infrastructure: %v\n",err)
            return 1
        }
    }
    if service=="integration" || service=="all" {
        if err:=integration.Migrate(ctx,db);err!=nil{
            _,_=fmt.Fprintf(os.Stderr,"migrate integration: %v\n",err)
            return 1
        }
    }
    if service=="identity" || service=="all" {
        if err:=identity.Migrate(ctx,db);err!=nil{
            _,_=fmt.Fprintf(os.Stderr,"migrate identity: %v\n",err)
            return 1
        }
    }
    if service=="accounting" || service=="all" {
        if err:=accounting.Migrate(ctx,db);err!=nil{
            _,_=fmt.Fprintf(os.Stderr,"migrate accounting: %v\n",err)
            return 1
        }
    }
    if service=="all" { return migrateAudit(ctx,db) }
    return 0
}

func migrateAudit(ctx context.Context,shared *gorm.DB) int {
	db,closeDB,err:=bootstrap.NewAuditPostgresDB(ctx,shared)
	if err!=nil { _,_=fmt.Fprintln(os.Stderr,"open Audit migration database:",err);return 1 }
	defer closeDB()
	if err:=audit.Migrate(ctx,db);err!=nil { _,_=fmt.Fprintln(os.Stderr,"migrate audit:",err);return 1 }
	return 0
}
