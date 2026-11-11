package validators_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestIdentityRegistrationValidation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Identity Registration Validation")
}
