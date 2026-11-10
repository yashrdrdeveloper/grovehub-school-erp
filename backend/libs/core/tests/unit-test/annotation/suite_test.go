package annotation_test

import (
    "testing"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

func TestAnnotationCatalog(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Annotation Phase 1 Contract Suite")
}
