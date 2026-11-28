package helpers_test

import (
	"grovehub/libs/core/http/internal/src/helpers"
	"log/slog"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("LevelForStatus", func() {
	It("maps successful and redirect responses to info, client errors to warn, and 404 or server errors to error", func() {
		for status, expected := range map[int]slog.Level{200: slog.LevelInfo, 302: slog.LevelInfo, 400: slog.LevelWarn, 403: slog.LevelWarn, 404: slog.LevelError, 500: slog.LevelError, 503: slog.LevelError} {
			Expect(helpers.LevelForStatus(status)).To(Equal(expected))
		}
	})
})
