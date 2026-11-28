// Command gormtable invokes the annotation generator.
package main

import (
	"fmt"
	"os"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/gormtable"
)

func main() {
 if len(os.Args)!=3 {fmt.Fprintln(os.Stderr,"usage: gormtable <entity-directory> <generated-output-file>");os.Exit(2)}
 if err:=gormtable.GenerateInto(os.Args[1],os.Args[2]);err!=nil {fmt.Fprintln(os.Stderr,err);os.Exit(1)}
}
