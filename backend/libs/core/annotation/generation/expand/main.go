// Command expand invokes the annotation generator.
package main

import (
	"fmt"
	"os"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/expand"
)

func main() {
 if len(os.Args)!=2 { fmt.Fprintln(os.Stderr,"usage: expand <staging-entity-directory>"); os.Exit(2) }
 if err:=expand.Expand(os.Args[1]);err!=nil {fmt.Fprintln(os.Stderr,err);os.Exit(1)}
}
