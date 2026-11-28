package logging_test

import (
	"log/slog"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corelogging "grovehub/libs/core/logging"
)

var _ = Describe("Logging facade", func() {
	It("exposes default and explicitly constructed levels", func() {
		defaults := corelogging.DefaultLevels()
		Expect(defaults.Enabled(slog.LevelInfo)).To(BeTrue())
		Expect(defaults.Enabled(slog.LevelDebug)).To(BeFalse())
		explicit := corelogging.NewLevels(slog.LevelDebug)
		Expect(explicit.Enabled(slog.LevelDebug)).To(BeTrue())
		Expect(explicit.Enabled(slog.LevelInfo)).To(BeFalse())
	})
})
