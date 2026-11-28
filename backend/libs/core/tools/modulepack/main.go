// Command modulepack prepares a Go module proxy release candidate.
package main

import (
	"fmt"
	"os"

	"grovehub/libs/core/tools/modulepack/internal/src/operations/pack"
)

func main() {
	if err := pack.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
