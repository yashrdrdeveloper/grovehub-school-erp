package helpers_test

import (
	"log/slog"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/logging/internal/src/helpers"
)

var _ = Describe("StandardLevel", func() {
	It("maps each supported name to its slog level", func() {
		for name, expected := range map[string]slog.Level{"DEBUG": slog.LevelDebug, "INFO": slog.LevelInfo, "WARN": slog.LevelWarn, "ERROR": slog.LevelError} {
			level, ok := helpers.StandardLevel(name)
			Expect(ok).To(BeTrue())
			Expect(level).To(Equal(expected))
		}
	})

	It("rejects unsupported names", func() {
		_, ok := helpers.StandardLevel("TRACE")
		Expect(ok).To(BeFalse())
	})
})
