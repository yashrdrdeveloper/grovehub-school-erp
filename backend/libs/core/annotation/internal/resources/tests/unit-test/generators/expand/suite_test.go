package expand_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestExpand(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Expand Suite")
}
