package validation_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAnnotationValidation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DTO Tag Normalization and Validation")
}
