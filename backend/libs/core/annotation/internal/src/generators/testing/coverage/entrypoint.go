package coverage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	generatedoutput "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/helpers/output"
)

// CoverEntrypoint instruments the actual executable and merges its runtime
// coverage into the unit-test profile. A safe CLI invocation exercises main
// without starting servers, opening production databases, or bypassing os.Exit.
func CoverEntrypoint(root, profile, pattern string, args []string, expectedExit int) error {
	directory, err := os.MkdirTemp(filepath.Dir(profile), "entrypoint-coverage-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(directory)

	binary := filepath.Join(directory, "command")
	if err := command(context.Background(), root, "go", "build", "-cover", "-covermode=atomic", "-coverpkg="+pattern, "-o", binary, pattern); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	run := exec.CommandContext(ctx, binary, args...)
	run.Dir = root
	run.Env = append(os.Environ(), "GOCOVERDIR="+directory)
	run.Stdout, run.Stderr = os.Stdout, os.Stderr

	exitCode := 0
	if err := run.Run(); err != nil {
		var status *exec.ExitError
		if !errors.As(err, &status) {
			return err
		}

		exitCode = status.ExitCode()
	}
	if exitCode != expectedExit {
		return fmt.Errorf("entrypoint exited with %d; expected %d", exitCode, expectedExit)
	}

	entrypointProfile := filepath.Join(directory, "entrypoint.out")
	if err := command(context.Background(), root, "go", "tool", "covdata", "textfmt", "-i="+directory, "-o="+entrypointProfile); err != nil {
		return err
	}

	return MergeEntrypointProfile(profile, entrypointProfile)
}

// MergeEntrypointProfile preserves every recorded block. Score subsequently
// combines duplicate block counters and checks the exact statement denominator.
func MergeEntrypointProfile(profile, entrypointProfile string) error {
	unit, err := os.ReadFile(profile)
	if err != nil {
		return err
	}

	entrypoint, err := os.ReadFile(entrypointProfile)
	if err != nil {
		return err
	}

	unitMode, _, _ := strings.Cut(string(unit), "\n")
	entrypointMode, records, _ := strings.Cut(string(entrypoint), "\n")
	if unitMode != "mode: atomic" || entrypointMode != unitMode || strings.TrimSpace(records) == "" {
		return errors.New("entrypoint coverage requires matching atomic profiles and recorded statements")
	}

	combined := strings.TrimRight(string(unit), "\n") + "\n" + records
	return generatedoutput.WriteAtomic(profile, "entrypoint-merged-*.out", []byte(combined))
}
