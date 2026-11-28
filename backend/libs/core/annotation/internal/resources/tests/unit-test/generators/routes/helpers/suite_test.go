package helpers_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRouteHelpers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RouteHelpers Suite")
}
