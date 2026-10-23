package logging_test

import (
	"log/slog"
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type levelCase struct {
	name       string
	statusCode int
	expected   slog.Level
}

var levelCases = []levelCase{
	{name: "success", statusCode: http.StatusOK, expected: slog.LevelInfo},
	{name: "redirect", statusCode: http.StatusFound, expected: slog.LevelInfo},
	{name: "bad request", statusCode: http.StatusBadRequest, expected: slog.LevelWarn},
	{name: "not found", statusCode: http.StatusNotFound, expected: slog.LevelError},
	{name: "server error", statusCode: http.StatusInternalServerError, expected: slog.LevelError},
}

func TestCoreHTTPLoggingUnit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Core HTTP Logging Unit Tests")
}
