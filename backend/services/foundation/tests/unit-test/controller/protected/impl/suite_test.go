package impl_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFoundationPersonController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Foundation Person Controller")
}
