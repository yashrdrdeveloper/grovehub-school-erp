package person_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPersonDTOValidation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Person DTO Tag Normalization and Validation")
}
