package observability_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestObservabilityCatalog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ObservabilityCatalog Suite")
}
