package errors_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCoreErrors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Core Error Definitions")
}
