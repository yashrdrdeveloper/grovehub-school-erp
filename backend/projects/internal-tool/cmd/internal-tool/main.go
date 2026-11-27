package main

import (
 "fmt"
 "log"
 "os"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/operatoradmin"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/launcher"
)

func main() {
 if len(os.Args)>1{
  if err:=operatoradmin.Run(os.Args[1:],os.Stdin);err!=nil{log.Fatal(err)}
  fmt.Println("Operator administration command completed")
  return
 }
 if err:=launcher.Run();err!=nil{log.Fatal(err)}
}
