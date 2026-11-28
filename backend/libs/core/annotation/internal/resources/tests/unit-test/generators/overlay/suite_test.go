package overlay_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestOverlay(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Annotation overlay Unit Tests")
}
