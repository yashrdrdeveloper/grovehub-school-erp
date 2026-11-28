package impl_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUnitOfWorkImplementation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UnitOfWorkImplementation Suite")
}
