package mappers_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMappers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Annotation mappers Unit Tests")
}
