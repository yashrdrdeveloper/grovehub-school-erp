package wrap_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestWrapOperation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "WrapOperation Suite")
}
