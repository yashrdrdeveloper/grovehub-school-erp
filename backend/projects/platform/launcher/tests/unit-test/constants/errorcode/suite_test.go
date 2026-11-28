package errorcode_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	expectedRouteNotFoundCode    = "PLATFORM_ROUTE_NOT_FOUND"
	expectedRouteNotFoundMessage = "API route not found"
)

func TestLauncherErrorCodeUnit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Launcher Error Code Unit Tests")
}
