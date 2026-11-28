package testgen_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTestgen(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Annotation testgen Unit Tests")
}
