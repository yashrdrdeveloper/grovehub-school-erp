package contracts_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEventContracts(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EventContracts Suite")
}
