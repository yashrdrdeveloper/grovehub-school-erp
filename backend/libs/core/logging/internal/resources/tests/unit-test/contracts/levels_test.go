package contracts_test

import (
	"log/slog"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/logging/internal/src/contracts"
)

var _ = Describe("NewLevels", func() {
	It("retains exactly the selected levels and copies the input", func() {
		input := []slog.Level{slog.LevelError, slog.LevelDebug, slog.LevelError}
		levels := contracts.NewLevels(input...)
		input[0] = slog.LevelInfo
		Expect(levels.Enabled(slog.LevelError)).To(BeTrue())
		Expect(levels.Enabled(slog.LevelDebug)).To(BeTrue())
		Expect(levels.Enabled(slog.LevelInfo)).To(BeFalse())
		Expect(levels.Enabled(slog.LevelWarn)).To(BeFalse())
	})

	It("supports empty and zero-value sets without enabling any level", func() {
		Expect(contracts.NewLevels().Enabled(slog.LevelInfo)).To(BeFalse())
		Expect(contracts.Levels{}.Enabled(slog.LevelError)).To(BeFalse())
	})
})
