package errorcode_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	expectedServiceDownCode    = "IDENTITY_SERVICE_DOWN"
	expectedServiceDownMessage = "Identity service is unavailable"
)

func TestIdentityErrorCodeUnit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Identity Error Code Unit Tests")
}
