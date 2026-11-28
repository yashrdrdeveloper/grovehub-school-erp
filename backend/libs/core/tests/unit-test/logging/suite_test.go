package logging_test

import (
	"bytes"
	"log/slog"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corelogging "grovehub/libs/core/logging"
)

const (
	errorWarnLevels = "ERROR,WARN"
	allLevels       = "DEBUG,INFO,WARN,ERROR"
)

func TestCoreLoggingUnit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Core Logging Unit Tests")
}

func newFilteredLogger(levelsValue string) (*slog.Logger, *bytes.Buffer, error) {
	levels, err := corelogging.ParseLevels(levelsValue)
	if err != nil {
		return nil, nil, err
	}

	buffer := &bytes.Buffer{}
	handler := slog.NewJSONHandler(buffer, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(corelogging.NewLevelFilterHandler(handler, levels))
	return logger, buffer, nil
}
