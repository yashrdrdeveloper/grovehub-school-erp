package external_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestExternalCatalog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ExternalCatalog Suite")
}
