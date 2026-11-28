package correlation_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCorrelationOperation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CorrelationOperation Suite")
}
