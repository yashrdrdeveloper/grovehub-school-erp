package main

import (
	"os"

	"grovehub/cmd/internal/src/bootstrap/impl"
)

func main() {
	os.Exit(impl.New().Build().Run(os.Args[1:], os.Stdout, os.Stderr))
}
