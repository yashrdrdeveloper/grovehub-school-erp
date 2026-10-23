package bootstrap_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/launcher/bootstrap"
)

var _ = Describe("Bootstrap logger", func() {
	It("uses the exact levels configured in GROVEHUB_LOG_LEVELS", func() {
		setLogLevels("ERROR,WARN")
		buffer := newLogBuffer()

		logger, err := bootstrap.NewLogger(buffer)
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

	It("returns an error for an unsupported configured level", func() {
		setLogLevels("ERROR,TRACE")

		logger, err := bootstrap.NewLogger(newLogBuffer())
		Expect(logger).To(BeNil())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("unsupported log level"))
	})
})
