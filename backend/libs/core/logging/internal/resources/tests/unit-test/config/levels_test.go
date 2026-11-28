package config_test

import (
	"log/slog"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/logging/internal/src/config"
)

var _ = Describe("DefaultLevels", func() {
	It("enables INFO, WARN and ERROR only", func() {
		levels := config.DefaultLevels()
		Expect(levels.Enabled(slog.LevelDebug)).To(BeFalse())
		Expect(levels.Enabled(slog.LevelInfo)).To(BeTrue())
		Expect(levels.Enabled(slog.LevelWarn)).To(BeTrue())
		Expect(levels.Enabled(slog.LevelError)).To(BeTrue())
	})
})

var _ = Describe("ParseLevels", func() {
	It("uses defaults for empty or whitespace-only configuration", func() {
		for _, value := range []string{"", " \t"} {
			levels, err := config.ParseLevels(value)
			Expect(err).NotTo(HaveOccurred())
			Expect(levels).To(Equal(config.DefaultLevels()))
		}
	})

	It("normalizes case, trims whitespace, skips empty entries and deduplicates", func() {
		levels, err := config.ParseLevels(" error, ,warn, ERROR, ")
		Expect(err).NotTo(HaveOccurred())
		Expect(levels.Enabled(slog.LevelError)).To(BeTrue())
		Expect(levels.Enabled(slog.LevelWarn)).To(BeTrue())
		Expect(levels.Enabled(slog.LevelInfo)).To(BeFalse())
		Expect(levels.Enabled(slog.LevelDebug)).To(BeFalse())
	})

	It("rejects unknown levels and delimiter-only configuration", func() {
		levels, err := config.ParseLevels("INFO,TRACE")
		Expect(err).To(MatchError(ContainSubstring("unsupported log level")))
		Expect(levels.Enabled(slog.LevelInfo)).To(BeFalse())
		levels, err = config.ParseLevels(", ,")
		Expect(err).To(MatchError("at least one log level is required"))
		Expect(levels.Enabled(slog.LevelError)).To(BeFalse())
	})
})
