// Package registry is the School-owned school directory and lifecycle boundary.
package registry

import (
 "context"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/registry/internal/app/services/delivery"
 "database/sql"
 "net/http"

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
 mux.HandleFunc(restricted.Pattern,handler)
 mux.HandleFunc(restricted.RenamePattern,rename)
 return nil
}

type DeliveryConfig=delivery.Config
func StartDelivery(ctx context.Context,db *sql.DB,c DeliveryConfig)(func(),error){
 publisher,err:=delivery.Client(c);if err!=nil{return nil,err}
 if publisher==nil{return func(){},nil}
 run,stop:=context.WithCancel(ctx);done:=make(chan struct{})
 go func(){defer close(done);delivery.New(db,c.SchoolID,publisher).Run(run)}()
 return func(){stop();<-done},nil
}
