// Package execution exposes the shared command execution contract to CMD.
package execution

import "grovehub/libs/platform/command"

// Executor forwards process input and output without knowing individual commands.
type Executor = command.Command
