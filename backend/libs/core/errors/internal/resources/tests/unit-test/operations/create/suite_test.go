package create_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCreateOperation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CreateOperation Suite")
}
