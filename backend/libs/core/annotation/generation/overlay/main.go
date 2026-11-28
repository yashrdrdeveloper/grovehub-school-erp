// Command overlay invokes the annotation generator.
package main

import (
	"fmt"
	"os"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/overlay"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: overlay <application-directory> <core-directory>")
		os.Exit(2)
	}
	if err := overlay.WriteOverlay(os.Args[1], os.Args[2]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
