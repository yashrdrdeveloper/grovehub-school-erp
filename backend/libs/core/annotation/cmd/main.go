// Command annotation runs the Core annotation processors.
package main

import (
	"os"

	"grovehub/libs/core/annotation/internal/src/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
