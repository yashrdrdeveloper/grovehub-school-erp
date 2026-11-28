package main

import (
	"os"

	"grovehub/cmd/internal/src/bootstrap"
)

func main() {
	os.Exit(bootstrap.New().Run(os.Args[1:], os.Stdout, os.Stderr))
}
