package module_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCoreModule(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Core Module Suite")
}
