// Command modulepack prepares a Go module proxy release candidate.
package main

import (
	"fmt"
	"os"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/tools/modulepack/internal/src/operations/pack"
)

func main() {
	if err := pack.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
