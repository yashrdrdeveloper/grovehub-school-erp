package validators_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestProjectConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ProjectConfig Suite")
}
