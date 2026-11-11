package request_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPersonPatchRequest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Foundation Person PATCH request")
}
