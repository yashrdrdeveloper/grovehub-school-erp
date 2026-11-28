package main

import (
	"os"

	"grovehub/launcher"
)

func main() {
	os.Exit(launcher.Execute(os.Args[1:], os.Stdout, os.Stderr))
}
