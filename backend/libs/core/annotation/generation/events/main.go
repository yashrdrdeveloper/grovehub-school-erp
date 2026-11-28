// Command events invokes the annotation generator.
package main

import (
	"fmt"
	"os"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/events"
)

func main(){
    if len(os.Args)!=2{fmt.Fprintln(os.Stderr,"usage: events <backend-dir>");os.Exit(2)}
    if err:=events.GenerateAll(os.Args[1]);err!=nil{fmt.Fprintln(os.Stderr,err);os.Exit(1)}
}
