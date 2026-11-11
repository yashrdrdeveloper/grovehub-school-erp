package helpers_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRegistrationErrorResponse(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Identity Controller Errors")
}
