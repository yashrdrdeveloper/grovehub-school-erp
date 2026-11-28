package composition_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCompositionCatalog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CompositionCatalog Suite")
}
