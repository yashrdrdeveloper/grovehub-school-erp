package mapping_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMappingCatalog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MappingCatalog Suite")
}
