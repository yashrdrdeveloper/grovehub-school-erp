package catalogtesting_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTestingCatalog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TestingCatalog Suite")
}
