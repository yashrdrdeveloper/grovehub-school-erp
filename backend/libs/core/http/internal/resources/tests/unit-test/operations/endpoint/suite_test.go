package endpoint_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEndpoint(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Core HTTP Endpoint Suite")
}
