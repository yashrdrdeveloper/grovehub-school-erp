package main

import (
    "context"
    "fmt"
    "os"
    "time"

    "grovehub/launcher/bootstrap"
    "grovehub/services/foundation"
    "grovehub/services/infrastructure"
)

// runMigrationCommand applies embedded service migrations before app startup.
func runMigrationCommand() int {
    if len(os.Args)!=3 || os.Args[1]!="migrate" {
        _,_=fmt.Fprintln(os.Stderr,"usage: groovehub-platform migrate foundation|infrastructure|all")
        return 2
    }
    service:=os.Args[2]
    if service!="foundation" && service!="infrastructure" && service!="all" {
        _,_=fmt.Fprintln(os.Stderr,"usage: groovehub-platform migrate foundation|infrastructure|all")
        return 2
    }
    ctx,cancel:=context.WithTimeout(context.Background(),10*time.Minute)
    defer cancel()
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
    return 0
}
