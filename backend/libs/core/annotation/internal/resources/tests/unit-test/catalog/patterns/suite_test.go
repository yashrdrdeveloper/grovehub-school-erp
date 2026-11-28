package patterns_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPatternsCatalog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PatternsCatalog Suite")
}
