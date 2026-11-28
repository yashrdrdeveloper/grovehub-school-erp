package strings_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestStringOperations(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "StringOperations Suite")
}
