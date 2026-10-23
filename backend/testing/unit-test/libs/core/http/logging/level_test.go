package logging_test

import (
	"log/slog"
	"net/http"
	"testing"

	httplogging "grovehub/libs/core/http/logging"
)

func TestLevelForStatus(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		expected   slog.Level
	}{
		{name: "success", statusCode: http.StatusOK, expected: slog.LevelInfo},
		{name: "redirect", statusCode: http.StatusFound, expected: slog.LevelInfo},
		{name: "bad request", statusCode: http.StatusBadRequest, expected: slog.LevelWarn},
		{name: "not found", statusCode: http.StatusNotFound, expected: slog.LevelError},
		{name: "server error", statusCode: http.StatusInternalServerError, expected: slog.LevelError},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := httplogging.LevelForStatus(test.statusCode)
			if actual != test.expected {
				t.Fatalf("expected level %s, got %s", test.expected, actual)
			}
		})
	}
}
