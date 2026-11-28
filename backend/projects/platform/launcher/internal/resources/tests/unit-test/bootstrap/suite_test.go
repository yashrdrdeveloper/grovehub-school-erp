package bootstrap_test

import (
	"bytes"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const logLevelsEnvironmentVariable = "GROVEHUB_LOG_LEVELS"

func setLogLevels(value string) {
	previousValue, existed := os.LookupEnv(logLevelsEnvironmentVariable)
	Expect(os.Setenv(logLevelsEnvironmentVariable, value)).To(Succeed())

	DeferCleanup(func() {
		if existed {
			Expect(os.Setenv(logLevelsEnvironmentVariable, previousValue)).To(Succeed())
			return
		}
		Expect(os.Unsetenv(logLevelsEnvironmentVariable)).To(Succeed())
	})
}

func newLogBuffer() *bytes.Buffer {
	return &bytes.Buffer{}
}
