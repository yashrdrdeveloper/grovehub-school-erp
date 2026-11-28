package command

import platform "grovehub/libs/platform/command"

// Command is the common operator strategy contract.
type Command = platform.Command
type Database[D any] = platform.Database[D]
