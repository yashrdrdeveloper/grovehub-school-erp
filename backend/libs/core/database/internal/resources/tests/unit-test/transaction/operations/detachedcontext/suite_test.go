package detachedcontext_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDetachedContext(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DetachedContext Suite")
}
