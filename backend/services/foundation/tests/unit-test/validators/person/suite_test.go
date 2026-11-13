package person_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFoundationPersonValidation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Foundation Person Validation")
}
