// Command operator-admin explicitly migrates the private database or bootstraps
// its first L2 operator. The password is read from stdin, never argv.
package main

import (
 "context"
 "fmt"
 "io"
 "log"
 "os"
 "strings"
 "time"

 "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/postgres"
 auth "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/services/authentication"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/resources/migrations"
)
func run()error{
 if len(os.Args)<2{return fmt.Errorf("usage: operator-admin migrate | bootstrap EMAIL (password on stdin)")}
 ctx,cancel:=context.WithTimeout(context.Background(),30*time.Second);defer cancel()
 db,err:=postgres.OpenProjectPool(ctx,os.Getenv("INTERNAL_TOOL_DATABASE_URL"));if err!=nil{return err};defer db.Close()
 switch os.Args[1]{
 case "migrate":
  if len(os.Args)!=2{return fmt.Errorf("usage: operator-admin migrate")}
  return migrations.Apply(ctx,db)
 case "bootstrap":
  if len(os.Args)!=3{return fmt.Errorf("usage: operator-admin bootstrap EMAIL")}
  if err:=migrations.Check(ctx,db);err!=nil{return err}
  raw,err:=io.ReadAll(io.LimitReader(os.Stdin,259));if err!=nil||len(raw)>258{return fmt.Errorf("invalid password input")}
  password:=strings.TrimSuffix(strings.TrimSuffix(string(raw),"\n"),"\r")
  return auth.New(db).Bootstrap(ctx,os.Args[2],password)
 default:return fmt.Errorf("unknown operator-admin command")
 }
}
func main(){if err:=run();err!=nil{log.Fatal(err)};fmt.Println("Operator administration command completed")}
