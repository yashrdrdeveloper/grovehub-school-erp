// Package registry is the School-owned school directory and lifecycle boundary.
package registry

import (
 "context"
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
