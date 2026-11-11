package helpers_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPersonErrorResponse(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Foundation Controller Errors")
}
