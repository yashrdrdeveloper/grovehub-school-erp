package cli

import "io"

// command defines one executable processor and its argument contract.
type command struct {
	usage string
	run   func([]string, io.Writer) error
	min   int
	max   int
}

func oneArgument(usage string, execute func(string) error) command {
	return command{
		usage: usage,
		min:   1,
		max:   1,
		run: func(args []string, _ io.Writer) error {
			return execute(args[0])
		},
	}
}

func twoArguments(usage string, execute func(string, string) error) command {
	return command{
		usage: usage,
		min:   2,
		max:   2,
		run: func(args []string, _ io.Writer) error {
			return execute(args[0], args[1])
		},
	}
}
