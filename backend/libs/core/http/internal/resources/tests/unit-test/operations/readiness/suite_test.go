package readiness_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestReadiness(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Core HTTP Readiness Suite")
}
