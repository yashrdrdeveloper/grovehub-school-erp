package main

import (
 "fmt"
 "log"
 "os"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/operatoradmin"
)

func main(){
 if err:=operatoradmin.Run(os.Args[1:],os.Stdin);err!=nil{log.Fatal(err)}
 fmt.Println("Operator administration command completed")
}
