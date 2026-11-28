package constants_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPoolConstants(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PoolConstants Suite")
}
