package normalize_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCoreNormalization(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Core Normalization")
}
