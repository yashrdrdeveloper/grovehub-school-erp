package endpoint_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEndpointClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Typed HTTP Endpoint Client")
}
