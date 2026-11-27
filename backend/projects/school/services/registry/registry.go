// Package registry is the School-owned school directory and lifecycle boundary.
package registry

import (
 "context"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/registry/internal/app/services/delivery"
 "database/sql"
 "net/http"
 "log/slog"
 "sync"
 "time"

 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/registry/internal/app/controller/restricted"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/registry/internal/app/services"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/registry/internal/resources/migrations"
)
func Migrate(ctx context.Context,db *sql.DB)error{return migrations.Apply(ctx,db)}
func Check(ctx context.Context,db *sql.DB)error{return migrations.Check(ctx,db)}
func RegisterSchool(ctx context.Context,db *sql.DB,id,code,name string)error{return services.Register(ctx,db,id,code,name)}
func RegisterHTTP(mux *http.ServeMux,db *sql.DB,grants string,renameGrants ...string)error{
 handler,err:=restricted.New(db,grants);if err!=nil{return err}
 raw:="";if len(renameGrants)>0{raw=renameGrants[0]}
 rename,err:=restricted.NewRename(db,raw,grants);if err!=nil{return err}
 statusRaw:="";if len(renameGrants)>1{statusRaw=renameGrants[1]}
 status,err:=restricted.NewStatus(db,statusRaw,grants,raw);if err!=nil{return err}
 mux.HandleFunc(restricted.StatusPattern,status)
 mux.HandleFunc(restricted.Pattern,handler)
 mux.HandleFunc(restricted.RenamePattern,rename)
 return nil
}

type DeliveryConfig=delivery.Config
func ConfigureDelivery(raw string,legacy DeliveryConfig,workloadID,workloadPrivate,readRaw,renameRaw string,statusRaw ...string)([]DeliveryConfig,error){
 return delivery.Configure(raw,legacy,workloadID,workloadPrivate,readRaw,renameRaw,statusRaw...)
}
func DeliveryHealth(db *sql.DB,configs []DeliveryConfig)http.HandlerFunc{return delivery.Health(db,configs)}

// StartDeliveries builds every client before launching a bounded set of
// workers. A worker has a dedicated school credential and cannot select
// another school's outbox rows.
func StartDeliveries(ctx context.Context,db *sql.DB,configs []DeliveryConfig)(func(),error){
 workers:=make([]interface{Step(context.Context)(bool,error)},0,len(configs))
 for _,c:=range configs{
  if c.Operation=="status"{
   publisher,err:=delivery.StatusClient(c);if err!=nil{return nil,err}
   workers=append(workers,delivery.NewStatus(db,c.SchoolID,publisher));continue
  }
  publisher,err:=delivery.Client(c);if err!=nil{return nil,err}
  if publisher!=nil{workers=append(workers,delivery.New(db,c.SchoolID,publisher))}
 }
 run,stop:=context.WithCancel(ctx)
 var wait sync.WaitGroup
 // Round-robin shards bound database and network concurrency to four.
 n:=len(workers);if n>4{n=4}
 for shard:=0;shard<n;shard++{
  wait.Add(1)
  go func(start int){
   defer wait.Done()
   for run.Err()==nil{
    worked:=false
    for index:=start;index<len(workers)&&run.Err()==nil;index+=n{
     ok,err:=workers[index].Step(run)
     if err!=nil&&run.Err()==nil{slog.Warn("school receipt relay step failed")}
     if err==nil&&ok{worked=true}
    }
    if worked{continue}
    select{case <-run.Done():return;case <-time.After(time.Second):}
   }
  }(shard)
 }
 return func(){stop();wait.Wait()},nil
}
func StartDelivery(ctx context.Context,db *sql.DB,c DeliveryConfig)(func(),error){
 return StartDeliveries(ctx,db,[]DeliveryConfig{c})
}
