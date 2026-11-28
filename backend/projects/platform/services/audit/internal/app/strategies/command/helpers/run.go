package helpers

import (
	"errors"
	"flag"
	"fmt"
	"io"

	"grovehub/services/audit/internal/app/strategies/command/constants"
)

func Run(diagnostics io.Writer, prefix string, execute func() error) int {
	err := execute()
	if err == nil || errors.Is(err, flag.ErrHelp) {
		return constants.Success
	}

	if prefix == "" {
		_, _ = fmt.Fprintln(diagnostics, err)
	} else {
		_, _ = fmt.Fprintln(diagnostics, prefix, err)
	}

	return constants.Failure
}
