package impl_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestIdentityAuthController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Identity Registration Controller")
}
