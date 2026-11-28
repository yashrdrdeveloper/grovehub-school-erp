package main

import (
 "context"
 "errors"
 "log"
 "os"
 "time"

 coredatabase "grovehub/libs/core/database"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/launcher"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/registry"
)
func run(args []string)error{
 if len(args)==0{return launcher.Run()}
 if !(len(args)==1&&args[0]=="migrate")&&!(len(args)==4&&args[0]=="register-school"){return errors.New("usage: school [migrate | register-school UUID CODE NAME]")}
 ctx,cancel:=context.WithTimeout(context.Background(),time.Minute);defer cancel()
 db,err:=coredatabase.OpenProjectPool(ctx,os.Getenv("SCHOOL_DATABASE_URL"));if err!=nil{return err};defer db.Close()
 if args[0]=="migrate"{return registry.Migrate(ctx,db)}
 if err=registry.Check(ctx,db);err!=nil{return err}
 return registry.RegisterSchool(ctx,db,args[1],args[2],args[3])
}
func main(){if err:=run(os.Args[1:]);err!=nil{log.Fatal(err)}}
