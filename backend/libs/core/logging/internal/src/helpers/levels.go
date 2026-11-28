package helpers

import (
	"log/slog"

	"grovehub/libs/core/logging/internal/src/constants"
)

func StandardLevel(name string) (slog.Level, bool) {
	switch name {
	case constants.Debug:
		return slog.LevelDebug, true
	case constants.Info:
		return slog.LevelInfo, true
	case constants.Warn:
		return slog.LevelWarn, true
	case constants.Error:
		return slog.LevelError, true
	default:
		return 0, false
	}
}
