package filter_test

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/logging/internal/src/contracts"
	"grovehub/libs/core/logging/internal/src/operations/filter"
)

type failingWriter struct{ err error }

func (writer failingWriter) Write(_ []byte) (int, error) {
	return 0, writer.err
}

var _ = Describe("NewLevelFilterHandler", func() {
	It("emits exact enabled levels without treating them as a minimum threshold", func() {
		var buffer bytes.Buffer
		handler := filter.NewLevelFilterHandler(slog.NewJSONHandler(&buffer, &slog.HandlerOptions{Level: slog.LevelDebug}), contracts.NewLevels(slog.LevelDebug, slog.LevelError))
		logger := slog.New(handler)
		logger.Debug("debug message")
		logger.Info("info message")
		logger.Warn("warn message")
		logger.Error("error message")
		Expect(buffer.String()).To(ContainSubstring("debug message"))
		Expect(buffer.String()).To(ContainSubstring("error message"))
		Expect(buffer.String()).NotTo(ContainSubstring("info message"))
		Expect(buffer.String()).NotTo(ContainSubstring("warn message"))
	})

	It("respects downstream level filtering as well as the selected levels", func() {
		handler := filter.NewLevelFilterHandler(slog.NewTextHandler(GinkgoWriter, &slog.HandlerOptions{Level: slog.LevelError}), contracts.NewLevels(slog.LevelDebug, slog.LevelError))
		Expect(handler.Enabled(context.Background(), slog.LevelDebug)).To(BeFalse())
		Expect(handler.Enabled(context.Background(), slog.LevelError)).To(BeTrue())
		Expect(handler.Enabled(context.Background(), slog.LevelInfo)).To(BeFalse())
	})

	It("retains filtering, attributes and groups on derived handlers", func() {
		var buffer bytes.Buffer
		handler := filter.NewLevelFilterHandler(slog.NewJSONHandler(&buffer, nil), contracts.NewLevels(slog.LevelWarn))
		derived := handler.WithAttrs([]slog.Attr{slog.String("service", "test")}).WithGroup("request")
		Expect(derived.Enabled(context.Background(), slog.LevelInfo)).To(BeFalse())
		logger := slog.New(derived)
		logger.Warn("accepted", "id", "one")
		Expect(buffer.String()).To(ContainSubstring(`"service":"test"`))
		Expect(buffer.String()).To(ContainSubstring(`"request":{"id":"one"}`))
	})

	It("returns downstream write errors from Handle", func() {
		failure := errors.New("write failed")
		handler := filter.NewLevelFilterHandler(slog.NewJSONHandler(failingWriter{err: failure}, nil), contracts.NewLevels(slog.LevelError))
		record := slog.NewRecord(time.Now(), slog.LevelError, "error", 0)
		Expect(handler.Handle(context.Background(), record)).To(MatchError(failure))
	})
})
