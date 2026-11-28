package unitofwork_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUnitOfWork(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UnitOfWork Suite")
}
