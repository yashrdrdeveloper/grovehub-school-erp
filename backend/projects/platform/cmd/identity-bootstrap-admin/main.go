package main

import (
	"fmt"
	"os"

	"grovehub/services/identity/bootstrapadmin"
)

func main() {
	if err := bootstrapadmin.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "bootstrap-admin:", err)
		os.Exit(1)
	}
}
