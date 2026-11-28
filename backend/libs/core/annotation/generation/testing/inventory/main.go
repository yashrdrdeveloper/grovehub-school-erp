// Command inventory invokes the annotation generator.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/testing/inventory"
)

func main() {
    if len(os.Args)<2||len(os.Args)>3 {fmt.Fprintln(os.Stderr,"usage: inventory <package-owner-directory> [--strict]");os.Exit(2)}
    report,err:=inventory.Inspect(os.Args[1])
    if err!=nil{fmt.Fprintln(os.Stderr,err);os.Exit(1)}
    output,err:=json.MarshalIndent(report,"","  ")
    if err!=nil{fmt.Fprintln(os.Stderr,err);os.Exit(1)}
    fmt.Println(string(output))
    if len(os.Args)==3&&os.Args[2]=="--strict"&&len(report.Missing)>0 {os.Exit(1)}
}
