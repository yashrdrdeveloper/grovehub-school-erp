package logging_test

import (
	"context"
	"log/slog"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corelogging "grovehub/libs/core/logging"
)

var _ = Describe("Configurable log levels", func() {
	It("emits only the exact levels selected in the environment value", func() {
		logger, buffer, err := newFilteredLogger(errorWarnLevels)
		Expect(err).NotTo(HaveOccurred())

		logger.Debug("debug message")
		logger.Info("info message")
		logger.Warn("warn message")
		logger.Error("error message")

		output := buffer.String()
		Expect(output).NotTo(ContainSubstring("debug message"))
		Expect(output).NotTo(ContainSubstring("info message"))
		Expect(output).To(ContainSubstring("warn message"))
		Expect(output).To(ContainSubstring("error message"))
	})

	It("accepts every supported standard level", func() {
		levels, err := corelogging.ParseLevels(allLevels)
		Expect(err).NotTo(HaveOccurred())

		Expect(levels.Enabled(slog.LevelDebug)).To(BeTrue())
		Expect(levels.Enabled(slog.LevelInfo)).To(BeTrue())
		Expect(levels.Enabled(slog.LevelWarn)).To(BeTrue())
		Expect(levels.Enabled(slog.LevelError)).To(BeTrue())
	})

	It("parses levels case-insensitively and trims whitespace", func() {
		levels, err := corelogging.ParseLevels(" error, warn ")
		Expect(err).NotTo(HaveOccurred())
		Expect(levels.Enabled(slog.LevelWarn)).To(BeTrue())
		Expect(levels.Enabled(slog.LevelError)).To(BeTrue())
		Expect(levels.Enabled(slog.LevelInfo)).To(BeFalse())
	})

	It("uses INFO WARN and ERROR when the setting is absent", func() {
		levels, err := corelogging.ParseLevels("")
		Expect(err).NotTo(HaveOccurred())
		Expect(levels.Enabled(slog.LevelDebug)).To(BeFalse())
		Expect(levels.Enabled(slog.LevelInfo)).To(BeTrue())
		Expect(levels.Enabled(slog.LevelWarn)).To(BeTrue())
		Expect(levels.Enabled(slog.LevelError)).To(BeTrue())
	})

	It("rejects unknown levels", func() {
		_, err := corelogging.ParseLevels("ERROR,TRACE")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("unsupported log level"))
	})

	It("preserves filtering when logger attributes and groups are added", func() {
		levels, err := corelogging.ParseLevels(errorWarnLevels)
		Expect(err).NotTo(HaveOccurred())

		handler := corelogging.NewLevelFilterHandler(
			slog.NewTextHandler(GinkgoWriter, &slog.HandlerOptions{Level: slog.LevelDebug}),
			levels,
		)
		Expect(handler.WithAttrs(nil).Enabled(context.Background(), slog.LevelInfo)).To(BeFalse())
		Expect(handler.WithGroup("request").Enabled(context.Background(), slog.LevelError)).To(BeTrue())
	})
})
