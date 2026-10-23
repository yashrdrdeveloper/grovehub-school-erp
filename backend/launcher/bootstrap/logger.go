package bootstrap

import (
	"io"
	"log/slog"
	"os"

	corelogging "grovehub/libs/core/logging"
)

const logLevelsEnvironmentVariable = "GROVEHUB_LOG_LEVELS"

func NewLogger(output io.Writer) (*slog.Logger, error) {
	levels, err := corelogging.ParseLevels(os.Getenv(logLevelsEnvironmentVariable))
	if err != nil {
		return nil, err
	}

	jsonHandler := slog.NewJSONHandler(output, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	return slog.New(corelogging.NewLevelFilterHandler(jsonHandler, levels)), nil
}
