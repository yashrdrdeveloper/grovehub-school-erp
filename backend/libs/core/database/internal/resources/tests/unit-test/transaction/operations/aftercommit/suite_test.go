package aftercommit_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAfterCommit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AfterCommit Suite")
}
