package patchdto_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPatchdto(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Annotation patchdto Unit Tests")
}
