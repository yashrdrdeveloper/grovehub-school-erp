// Package operatoradmin implements explicit database and operator administration.
package operatoradmin

import (
 "context"
 "fmt"
 "io"
 "os"
 "strings"
 "time"

 coredatabase "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"
 auth "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/services/authentication"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/resources/migrations"
)
func Run(args []string, input io.Reader)error{
 if len(args)<1{return fmt.Errorf("usage: operator-admin migrate | bootstrap EMAIL (password on stdin)")}
 if (args[0]=="migrate"&&len(args)!=1)||(args[0]=="bootstrap"&&len(args)!=2)||(args[0]!="migrate"&&args[0]!="bootstrap"){return fmt.Errorf("usage: migrate | bootstrap EMAIL (password on stdin)")}
 ctx,cancel:=context.WithTimeout(context.Background(),30*time.Second);defer cancel()
 db,err:=coredatabase.OpenProjectPool(ctx,os.Getenv("INTERNAL_TOOL_DATABASE_URL"));if err!=nil{return err};defer db.Close()
 switch args[0]{
 case "migrate":
  if len(args)!=1{return fmt.Errorf("usage: operator-admin migrate")}
  return migrations.Apply(ctx,db)
 case "bootstrap":
  if len(args)!=2{return fmt.Errorf("usage: operator-admin bootstrap EMAIL")}
  if err:=migrations.Check(ctx,db);err!=nil{return err}
  raw,err:=io.ReadAll(io.LimitReader(input,259));if err!=nil||len(raw)>258{return fmt.Errorf("invalid password input")}
  password:=strings.TrimSuffix(strings.TrimSuffix(string(raw),"\n"),"\r")
  return auth.New(db).Bootstrap(ctx,args[1],password)
 default:return fmt.Errorf("unknown operator-admin command")
 }
}
