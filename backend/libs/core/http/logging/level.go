package logging

import (
	"log/slog"
	"net/http"
)

func LevelForStatus(statusCode int) slog.Level {
	switch {
	case statusCode == http.StatusNotFound:
		return slog.LevelError
	case statusCode >= http.StatusInternalServerError:
		return slog.LevelError
	case statusCode >= http.StatusBadRequest:
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}
