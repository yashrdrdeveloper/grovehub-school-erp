package helpers_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSearchHelpers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SearchHelpers Suite")
}
