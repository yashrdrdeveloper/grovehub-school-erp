// Command check invokes the annotation generator.
package main

import (
	"fmt"
	"os"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/check"
)

func main(){
    if len(os.Args)!=2 {fmt.Fprintln(os.Stderr,"usage: check <backend-directory>");os.Exit(2)}
    if err:=check.CheckServices(os.Args[1]);err!=nil {fmt.Fprintln(os.Stderr,err);os.Exit(1)}
}
