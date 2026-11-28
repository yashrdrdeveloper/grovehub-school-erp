// Package cli dispatches annotation commands without exposing processor packages.
package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"sort"
)

// Run returns a process exit code; only the executable calls os.Exit.
func Run(args []string, stdout, stderr io.Writer) int {
	registry := commands()

	if len(args) == 0 {
		printUsage(stderr, registry)
		return 2
	}

	if isHelp(args[0]) {
		printUsage(stdout, registry)
		return 0
	}

	selected, exists := registry[args[0]]
	if !exists {
		fmt.Fprintf(stderr, "unknown annotation command %q\n", args[0])
		printUsage(stderr, registry)
		return 2
	}

	parameters := args[1:]
	if len(parameters) == 1 && isHelp(parameters[0]) && args[0] != "coverage" {
		printCommandUsage(stdout, args[0], selected)
		return 0
	}

	if len(parameters) < selected.min || (selected.max >= 0 && len(parameters) > selected.max) {
		printCommandUsage(stderr, args[0], selected)
		return 2
	}

	err := selected.run(parameters, stdout)
	if errors.Is(err, flag.ErrHelp) {
		return 0
	}
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	return 0
}

func isHelp(argument string) bool {
	switch argument {
	case "help", "-h", "--help":
		return true
	default:
		return false
	}
}

func printUsage(output io.Writer, registry map[string]command) {
	fmt.Fprintln(output, "Usage: annotation <command> [arguments]")
	fmt.Fprintln(output, "\nCommands:")

	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}

	sort.Strings(names)

	for _, name := range names {
		fmt.Fprintf(output, "  %-12s %s\n", name, registry[name].usage)
	}
}

func printCommandUsage(output io.Writer, name string, selected command) {
	fmt.Fprintf(output, "Usage: annotation %s %s\n", name, selected.usage)
}
