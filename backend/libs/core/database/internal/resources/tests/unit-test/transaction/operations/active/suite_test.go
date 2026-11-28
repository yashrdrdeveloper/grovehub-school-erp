package active_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestActive(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Active Suite")
}
