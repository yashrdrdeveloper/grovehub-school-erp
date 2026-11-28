package contracts_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestErrorContracts(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ErrorContracts Suite")
}
