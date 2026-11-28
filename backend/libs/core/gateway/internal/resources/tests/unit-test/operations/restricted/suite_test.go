package restricted_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRestricted(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gateway restricted Suite")
}
