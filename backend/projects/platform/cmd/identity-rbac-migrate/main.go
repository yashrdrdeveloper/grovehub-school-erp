package main

import (
	"fmt"
	"os"

	"grovehub/services/identity/rbacmigrationcli"
)

func main() {
	if err := rbacmigrationcli.Run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "identity-rbac-migrate:", err)
		os.Exit(1)
	}
}
