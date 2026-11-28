package cli

import (
	"encoding/json"
	"fmt"
	"io"

	"grovehub/libs/core/annotation/internal/src/generators/testing/inventory"
)

func runInventory(args []string, output io.Writer) error {
	if len(args) == 2 && args[1] != "--strict" {
		return fmt.Errorf("unknown inventory option %q", args[1])
	}

	report, err := inventory.Inspect(args[0])
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(report); err != nil {
		return err
	}

	if len(args) == 2 && len(report.Missing) > 0 {
		return fmt.Errorf("inventory has %d functions without mirrored tests", len(report.Missing))
	}

	return nil
}
