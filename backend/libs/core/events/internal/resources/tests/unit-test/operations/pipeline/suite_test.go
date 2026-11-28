package pipeline_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPipelineOperation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PipelineOperation Suite")
}
