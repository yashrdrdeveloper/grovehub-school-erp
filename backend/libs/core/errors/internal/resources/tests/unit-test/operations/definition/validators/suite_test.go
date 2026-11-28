package validators_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestErrorStatusValidator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ErrorStatusValidator Suite")
}
