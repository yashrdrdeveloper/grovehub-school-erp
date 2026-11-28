package definition_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDefinitionOperation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DefinitionOperation Suite")
}
